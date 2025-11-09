package service

import (
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/valueobject"
)

type BalanceValidator struct{}

func NewBalanceValidator() *BalanceValidator {
	return &BalanceValidator{}
}

// ValidateDeposit validates if a deposit can be made to a wallet
func (bv *BalanceValidator) ValidateDeposit(wallet *entity.Wallet, amount valueobject.Money) error {
	return wallet.CanDeposit(amount)
}

// GetMaxAllowedDeposit calculates the maximum amount that can be deposited
func (bv *BalanceValidator) GetMaxAllowedDeposit(wallet *entity.Wallet) (valueobject.Money, error) {
	maxBalance, err := wallet.Type.MaxBalance()
	if err != nil {
		return valueobject.Money{}, err
	}

	return maxBalance.Subtract(wallet.Balance)
}
