package repository

import (
	"context"
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/valueobject"
)

// WalletRepository defines the interface for wallet persistence
type WalletRepository interface {
	FindByAccountID(ctx context.Context, accountID valueobject.AccountID) (*entity.Wallet, error)
	FindByID(ctx context.Context, id int64) (*entity.Wallet, error)
	Create(ctx context.Context, wallet *entity.Wallet) error
	Update(ctx context.Context, wallet *entity.Wallet) error
	ExistsByAccountID(ctx context.Context, accountID valueobject.AccountID) (bool, error)
}
