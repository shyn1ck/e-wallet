package mapper

import (
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/valueobject"
	"e-wallet/internal/infrastructure/database/models"
	"fmt"
)

type TransactionMapper struct{}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (m *TransactionMapper) ToDomain(dbTx *models.Transaction) (*entity.Transaction, error) {
	amount, err := valueobject.NewMoneyFromMinor(dbTx.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	return &entity.Transaction{
		ID:        dbTx.ID,
		WalletID:  dbTx.WalletID,
		Type:      entity.TransactionType(dbTx.Type),
		Amount:    amount,
		CreatedAt: dbTx.CreatedAt,
	}, nil
}

func (m *TransactionMapper) ToModel(tx *entity.Transaction) *models.Transaction {
	return &models.Transaction{
		ID:        tx.ID,
		WalletID:  tx.WalletID,
		Type:      string(tx.Type),
		Amount:    tx.Amount.Amount(),
		CreatedAt: tx.CreatedAt,
	}
}
