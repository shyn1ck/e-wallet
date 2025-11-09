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

type WalletCheckUseCase struct {
	walletRepo repository.WalletRepository
}

func NewWalletCheckUseCase(walletRepo repository.WalletRepository) *WalletCheckUseCase {
	return &WalletCheckUseCase{
		walletRepo: walletRepo,
	}
}

// Execute checks if a wallet exists
func (uc *WalletCheckUseCase) Execute(ctx context.Context, req *request.CheckWalletRequest) (*response.CheckWalletResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrValidationFailed
	}

	accountID, err := valueobject.NewAccountID(req.AccountID)
	if err != nil {
		return nil, apperrors.ErrInvalidRequest
	}

	exists, err := uc.walletRepo.ExistsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	resp := &response.CheckWalletResponse{
		Exists: exists,
	}

	if exists {
		resp.AccountID = req.AccountID
	}

	return resp, nil
}
