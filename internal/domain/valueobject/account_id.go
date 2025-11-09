package valueobject

import (
	apperrors "e-wallet/pkg/errors"
	"regexp"
	"strings"
)

// AccountID represents a wallet account identifier
type AccountID struct {
	value string
}

// accountIDPattern is a regular expression for validating account IDs
var accountIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{3,50}$`)

// NewAccountID creates and validates an AccountID
func NewAccountID(value string) (AccountID, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return AccountID{}, apperrors.ErrEmptyAccountID
	}

	if !accountIDPattern.MatchString(trimmed) {
		return AccountID{}, apperrors.ErrInvalidAccountID
	}

	return AccountID{value: trimmed}, nil
}

func (a AccountID) Value() string {
	return a.value
}
func (a AccountID) String() string {
	return a.value
}

func (a AccountID) Equals(other AccountID) bool {
	return a.value == other.value
}
