package repository

import (
	"context"
	"e-wallet/internal/domain/entity"
	"time"
)

// TransactionRepository defines the interface for transaction persistence
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	FindByWalletID(ctx context.Context, walletID int64) ([]*entity.Transaction, error)
	GetMonthlyStats(ctx context.Context, walletID int64, month time.Time) (*MonthlyStats, error)
}

type MonthlyStats struct {
	TotalCount  int64
	TotalAmount int64 // (dirams)
}
