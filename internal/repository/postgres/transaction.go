package postgres

import (
	"context"
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/infrastructure/database"
	"e-wallet/internal/infrastructure/database/models"
	"e-wallet/internal/infrastructure/logger"
	"e-wallet/internal/repository/mapper"
	apperrors "e-wallet/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db     *gorm.DB
	mapper *mapper.TransactionMapper
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		mapper: mapper.NewTransactionMapper(),
	}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	db := database.GetDB(ctx, r.db)
	dbTx := r.mapper.ToModel(transaction)
	err := db.WithContext(ctx).Create(dbTx).Error
	if err != nil {
		logger.Error.Printf("[postgres.Create]: Failed to create transaction for wallet_id %d: %v", transaction.WalletID, err)
		return apperrors.TranslateError(err)
	}

	transaction.ID = dbTx.ID
	transaction.CreatedAt = dbTx.CreatedAt

	return nil
}

func (r *TransactionRepository) FindByWalletID(ctx context.Context, walletID int64) ([]*entity.Transaction, error) {
	db := database.GetDB(ctx, r.db)
	var dbTransactions []models.Transaction
	err := db.WithContext(ctx).Where("wallet_id = ?", walletID).Order("created_at DESC").Find(&dbTransactions).Error
	if err != nil {
		logger.Error.Printf("[postgres.FindByWalletID]: Failed to find transactions for wallet_id %d: %v", walletID, err)
		return nil, apperrors.TranslateError(err)
	}

	transactions := make([]*entity.Transaction, 0, len(dbTransactions))
	for _, dbTx := range dbTransactions {
		tx, err := r.mapper.ToDomain(&dbTx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// GetMonthlyStats retrieves monthly statistics for a wallet
func (r *TransactionRepository) GetMonthlyStats(ctx context.Context, walletID int64, month time.Time) (*repository.MonthlyStats, error) {
	firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)

	db := database.GetDB(ctx, r.db)
	var stats repository.MonthlyStats
	err := db.WithContext(ctx).
		Model(&models.Transaction{}).
		Select("COUNT(*) as total_count, COALESCE(SUM(amount), 0) as total_amount").
		Where("wallet_id = ? AND type = ? AND created_at BETWEEN ? AND ?",
			walletID, string(entity.TransactionTypeDeposit), firstDay, lastDay).
		Scan(&stats).Error

	if err != nil {
		logger.Error.Printf("[postgres.GetMonthlyStats]: Failed to get monthly stats for wallet_id %d: %v", walletID, err)
		return nil, apperrors.TranslateError(err)
	}

	return &stats, nil
}
