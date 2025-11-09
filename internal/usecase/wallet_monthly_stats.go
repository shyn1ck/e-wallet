package usecase

import (
	"context"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/domain/valueobject"
	"e-wallet/internal/dto/request"
	"e-wallet/internal/dto/response"
	apperrors "e-wallet/pkg/errors"
	"e-wallet/pkg/utils"
	"e-wallet/pkg/validator"
)

type WalletMonthlyStatsUseCase struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
}

func NewWalletMonthlyStatsUseCase(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
) *WalletMonthlyStatsUseCase {
	return &WalletMonthlyStatsUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

// Execute retrieves monthly statistics for a wallet
func (uc *WalletMonthlyStatsUseCase) Execute(ctx context.Context, req *request.GetMonthlyStatsRequest) (*response.MonthlyStatsResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrValidationFailed
	}

	accountID, err := valueobject.NewAccountID(req.AccountID)
	if err != nil {
		return nil, apperrors.ErrInvalidRequest
	}

	wallet, err := uc.walletRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	currentMonth := utils.GetCurrentMonth()

	stats, err := uc.transactionRepo.GetMonthlyStats(ctx, wallet.ID, currentMonth)
	if err != nil {
		return nil, err
	}

	return &response.MonthlyStatsResponse{
		AccountID:   accountID.Value(),
		Month:       currentMonth.Format("2006-01"),
		TotalCount:  stats.TotalCount,
		TotalAmount: stats.TotalAmount,
		Currency:    valueobject.CurrencyTJS,
	}, nil
}
