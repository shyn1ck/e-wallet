package entity

import (
	"e-wallet/internal/domain/valueobject"
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit TransactionType = "deposit"
)

type Transaction struct {
	ID        int64
	WalletID  int64
	Type      TransactionType
	Amount    valueobject.Money
	CreatedAt time.Time
}

func NewTransaction(walletID int64, txType TransactionType, amount valueobject.Money) *Transaction {
	return &Transaction{
		WalletID:  walletID,
		Type:      txType,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}
