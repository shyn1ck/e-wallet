package valueobject

import (
	apperrors "e-wallet/pkg/errors"
	"strings"
)

type WalletType string

const (
	WalletTypeIdentified   WalletType = "identified"
	WalletTypeUnidentified WalletType = "unidentified"
)

const (
	MaxBalanceIdentified   = 10000000 // 100,000 TJS in dirams
	MaxBalanceUnidentified = 1000000  // 10,000 TJS in dirams
)

func NewWalletType(value string) (WalletType, error) {
	normalized := WalletType(strings.ToLower(strings.TrimSpace(value)))

	switch normalized {
	case WalletTypeIdentified, WalletTypeUnidentified:
		return normalized, nil
	default:
		return "", apperrors.ErrInvalidWalletType
	}
}

func (wt WalletType) String() string {
	return string(wt)
}

func (wt WalletType) IsIdentified() bool {
	return wt == WalletTypeIdentified
}
func (wt WalletType) MaxBalance() (Money, error) {
	switch wt {
	case WalletTypeIdentified:
		return NewMoney(MaxBalanceIdentified)
	case WalletTypeUnidentified:
		return NewMoney(MaxBalanceUnidentified)
	default:
		return Money{}, apperrors.ErrInvalidWalletType
	}
}
