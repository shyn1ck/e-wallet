package postgres

import (
	"context"
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/valueobject"
	"e-wallet/internal/infrastructure/database/models"
	"e-wallet/internal/infrastructure/logger"
	"e-wallet/internal/repository/mapper"
	apperrors "e-wallet/pkg/errors"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db     *gorm.DB
	mapper *mapper.WalletMapper
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db:     db,
		mapper: mapper.NewWalletMapper(),
	}
}

// FindByAccountID retrieves a wallet by account ID
func (r *WalletRepository) FindByAccountID(ctx context.Context, accountID valueobject.AccountID) (*entity.Wallet, error) {
	var dbWallet models.Wallet
	err := r.db.WithContext(ctx).Where("account_id = ?", accountID.Value()).First(&dbWallet).Error
	if err != nil {
		logger.Error.Printf("[postgres.FindByAccountID]: Failed to find wallet by account_id %s: %v", accountID.Value(), err)
		return nil, apperrors.TranslateError(err)
	}

	return r.mapper.ToDomain(&dbWallet)
}

// FindByID retrieves a wallet by ID
func (r *WalletRepository) FindByID(ctx context.Context, id int64) (*entity.Wallet, error) {
	var dbWallet models.Wallet
	err := r.db.WithContext(ctx).First(&dbWallet, id).Error
	if err != nil {
		logger.Error.Printf("[postgres.FindByID]: Failed to find wallet by id %d: %v", id, err)
		return nil, apperrors.TranslateError(err)
	}

	return r.mapper.ToDomain(&dbWallet)
}

// Create creates a new wallet
func (r *WalletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	dbWallet := r.mapper.ToModel(wallet)
	err := r.db.WithContext(ctx).Create(dbWallet).Error
	if err != nil {
		logger.Error.Printf("[postgres.Create]: Failed to create wallet: %v", err)
		return apperrors.TranslateError(err)
	}

	wallet.ID = dbWallet.ID
	wallet.CreatedAt = dbWallet.CreatedAt
	wallet.UpdatedAt = dbWallet.UpdatedAt

	return nil
}

// Update updates an existing wallet
func (r *WalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	dbWallet := r.mapper.ToModel(wallet)
	err := r.db.WithContext(ctx).Save(dbWallet).Error
	if err != nil {
		logger.Error.Printf("[postgres.Update]: Failed to update wallet id %d: %v", wallet.ID, err)
		return apperrors.TranslateError(err)
	}
	return nil
}

// ExistsByAccountID checks if a wallet exists by account ID
func (r *WalletRepository) ExistsByAccountID(ctx context.Context, accountID valueobject.AccountID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Wallet{}).Where("account_id = ?", accountID.Value()).Count(&count).Error
	if err != nil {
		logger.Error.Printf("[postgres.ExistsByAccountID]: Failed to check wallet existence for account_id %s: %v", accountID.Value(), err)
		return false, apperrors.TranslateError(err)
	}
	return count > 0, nil
}
