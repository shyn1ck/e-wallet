package usecase

import (
	"context"
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/domain/service"
	"e-wallet/internal/domain/valueobject"
	"e-wallet/internal/dto/request"
	"e-wallet/internal/dto/response"
	"e-wallet/internal/infrastructure/database"
	"e-wallet/internal/infrastructure/logger"
	apperrors "e-wallet/pkg/errors"
	"e-wallet/pkg/validator"

	"gorm.io/gorm"
)

// WalletDepositUseCase handles wallet deposit operations
type WalletDepositUseCase struct {
	db               *gorm.DB
	walletRepo       repository.WalletRepository
	transactionRepo  repository.TransactionRepository
	balanceValidator *service.BalanceValidator
}

// NewWalletDepositUseCase creates a new WalletDepositUseCase
func NewWalletDepositUseCase(
	db *gorm.DB,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	balanceValidator *service.BalanceValidator,
) *WalletDepositUseCase {
	return &WalletDepositUseCase{
		db:               db,
		walletRepo:       walletRepo,
		transactionRepo:  transactionRepo,
		balanceValidator: balanceValidator,
	}
}

// Execute performs a deposit operation
func (uc *WalletDepositUseCase) Execute(ctx context.Context, req *request.DepositRequest) (*response.DepositResponse, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrValidationFailed
	}

	// Create value objects
	accountID, err := valueobject.NewAccountID(req.AccountID)
	if err != nil {
		return nil, apperrors.ErrInvalidRequest
	}

	amount, err := valueobject.NewMoney(req.Amount)
	if err != nil {
		return nil, apperrors.ErrInvalidAmount
	}

	// Start transaction
	var resp *response.DepositResponse
	err = uc.db.Transaction(func(tx *gorm.DB) error {
		txCtx := database.InjectTx(ctx, tx)

		wallet, err := uc.walletRepo.FindByAccountID(txCtx, accountID)
		if err != nil {
			return err
		}

		logger.Info.Printf("Depositing %d dirams to wallet %s (current balance: %d dirams)",
			amount.Dirams(), accountID.Value(), wallet.Balance.Dirams())

		// Validate deposit
		if err := uc.balanceValidator.ValidateDeposit(wallet, amount); err != nil {
			return err
		}

		// Perform deposit
		if err := wallet.Deposit(amount); err != nil {
			return err
		}

		// Update wallet
		if err := uc.walletRepo.Update(txCtx, wallet); err != nil {
			return err
		}

		// Create transaction record
		transaction := entity.NewTransaction(wallet.ID, entity.TransactionTypeDeposit, amount)
		if err := uc.transactionRepo.Create(txCtx, transaction); err != nil {
			return err
		}

		logger.Info.Printf("Deposit successful. New balance: %d dirams, Transaction ID: %d",
			wallet.Balance.Dirams(), transaction.ID)

		// Build response
		resp = &response.DepositResponse{
			Success:       true,
			AccountID:     accountID.Value(),
			Amount:        amount.Dirams(),
			NewBalance:    wallet.Balance.Dirams(),
			Currency:      valueobject.CurrencyTJS,
			TransactionID: transaction.ID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
