package entity

import (
	"e-wallet/internal/domain/valueobject"
	apperrors "e-wallet/pkg/errors"
	"time"
)

// Wallet represents a wallet account entity
type Wallet struct {
	ID        int64
	AccountID valueobject.AccountID
	Type      valueobject.WalletType
	Balance   valueobject.Money
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *Wallet) CanDeposit(amount valueobject.Money) error {
	newBalance := w.Balance.Add(amount)
	maxBalance, err := w.Type.MaxBalance()
	if err != nil {
		return err
	}

	if newBalance.IsGreaterThan(maxBalance) {
		return apperrors.ErrBalanceExceedsLimit
	}

	return nil
}

func (w *Wallet) Deposit(amount valueobject.Money) error {
	if err := w.CanDeposit(amount); err != nil {
		return err
	}

	w.Balance = w.Balance.Add(amount)
	w.UpdatedAt = time.Now()

	return nil
}
