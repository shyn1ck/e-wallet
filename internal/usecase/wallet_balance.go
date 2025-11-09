package usecase

import (
	"context"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/domain/valueobject"
	"e-wallet/internal/dto/request"
	"e-wallet/internal/dto/response"
	apperrors "e-wallet/pkg/errors"
	"e-wallet/pkg/validator"
)

// WalletBalanceUseCase handles wallet balance retrieval
type WalletBalanceUseCase struct {
	walletRepo repository.WalletRepository
}

// NewWalletBalanceUseCase creates a new WalletBalanceUseCase
func NewWalletBalanceUseCase(walletRepo repository.WalletRepository) *WalletBalanceUseCase {
	return &WalletBalanceUseCase{
		walletRepo: walletRepo,
	}
}

// Execute retrieves wallet balance
func (uc *WalletBalanceUseCase) Execute(ctx context.Context, req *request.GetBalanceRequest) (*response.GetBalanceResponse, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrValidationFailed
	}

	// Create account ID value object
	accountID, err := valueobject.NewAccountID(req.AccountID)
	if err != nil {
		return nil, apperrors.ErrInvalidRequest
	}

	// Find wallet
	wallet, err := uc.walletRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &response.GetBalanceResponse{
		AccountID: accountID.Value(),
		Balance:   wallet.Balance.Dirams(),
		Currency:  valueobject.CurrencyTJS,
	}, nil
}
