package valueobject

import (
	apperrors "e-wallet/pkg/errors"
	"fmt"
)

type Money struct {
	amount int64
}

func NewMoney(dirams int64) (Money, error) {
	if dirams < 0 {
		return Money{}, apperrors.ErrInvalidAmount
	}

	return Money{amount: dirams}, nil
}

func NewMoneyFromMinor(dirams int64) (Money, error) {
	if dirams < 0 {
		return Money{}, apperrors.ErrInvalidAmount
	}
	return Money{amount: dirams}, nil
}

// Amount returns in dirams
func (m Money) Amount() int64 {
	return m.amount
}

// Dirams returns the amount in dirams
func (m Money) Dirams() int64 {
	return m.amount
}

// Somoni returns the amount in major units (somoni)
func (m Money) Somoni() float64 {
	return float64(m.amount) / 100.0
}

// Add adds two Money values
func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.amount < other.amount {
		return Money{}, apperrors.ErrInsufficientFunds
	}
	return Money{amount: m.amount - other.amount}, nil
}

func (m Money) IsGreaterThan(other Money) bool {
	return m.amount > other.amount
}

func (m Money) IsLessThan(other Money) bool {
	return m.amount < other.amount
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f TJS", m.Somoni())
}
