package errors

import (
	"errors"
	"net/http"
)

type APIError struct {
	Code    string
	Message string
	Status  int
}

func (e *APIError) Error() string {
	return e.Code
}

var (
	ErrBalanceExceedsLimit   = &APIError{"balance_limit_exceeded", "Deposit would exceed wallet balance limit", http.StatusBadRequest}
	ErrInvalidAmount         = &APIError{"invalid_amount", "Invalid amount", http.StatusBadRequest}
	ErrInvalidSignature      = &APIError{"invalid_signature", "Invalid signature", http.StatusUnauthorized}
	ErrMissingAuthData       = &APIError{"missing_auth_data", "Missing authentication data", http.StatusUnauthorized}
	ErrClientNotFound        = &APIError{"invalid_credentials", "Invalid credentials", http.StatusUnauthorized}
	ErrClientInactive        = &APIError{"client_inactive", "Client is inactive", http.StatusUnauthorized}
	ErrRateLimitExceeded     = &APIError{"rate_limit_exceeded", "Too many requests, please try again later", http.StatusTooManyRequests}
	ErrInvalidRequest        = &APIError{"invalid_request", "Invalid request", http.StatusBadRequest}
	ErrValidationFailed      = &APIError{"validation_failed", "Validation failed", http.StatusBadRequest}
	ErrInternalServerError   = &APIError{"internal_server_error", "An unexpected error occurred", http.StatusInternalServerError}
	ErrInvalidData           = &APIError{"invalid_data", "Invalid data", http.StatusBadRequest}
	ErrDataTooLong           = &APIError{"data_too_long", "Data is too long", http.StatusBadRequest}
	ErrRequiredField         = &APIError{"required_field", "Required field", http.StatusBadRequest}
	ErrRelatedRecordNotFound = &APIError{"related_record_not_found", "Related record not found", http.StatusNotFound}
	ErrAlreadyExists         = &APIError{"already_exists", "Already exists", http.StatusBadRequest}
	ErrRecordNotFound        = &APIError{"record_not_found", "Record not found", http.StatusNotFound}
	ErrEmptyAccountID        = &APIError{"empty_account_id", "Account ID cannot be empty", http.StatusBadRequest}
	ErrInvalidAccountID      = &APIError{"invalid_account_id", "Invalid account ID", http.StatusBadRequest}
	ErrInsufficientFunds     = &APIError{"insufficient_funds", "Insufficient funds", http.StatusBadRequest}
	ErrInvalidWalletType     = &APIError{"invalid_wallet_type", "Invalid wallet type", http.StatusBadRequest}
)

// GetStatusCode returns HTTP status code
func GetStatusCode(err error) int {
	var e *APIError
	if errors.As(err, &e) {
		return e.Status
	}
	return http.StatusInternalServerError
}

// GetErrorMessage returns error message
func GetErrorMessage(err error) string {
	var e *APIError
	if errors.As(err, &e) {
		return e.Message
	}
	return "An unexpected error occurred"
}

// GetErrorCode returns error code
func GetErrorCode(err error) string {
	var e *APIError
	if errors.As(err, &e) {
		return e.Code
	}
	return "internal_server_error"
}
