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
	ErrBalanceExceedsLimit   = &APIError{"BALANCE_LIMIT_EXCEEDED", "Deposit would exceed wallet balance limit", http.StatusBadRequest}
	ErrInvalidAmount         = &APIError{"INVALID_AMOUNT", "Invalid amount", http.StatusBadRequest}
	ErrInvalidSignature      = &APIError{"INVALID_SIGNATURE", "Invalid HMAC signature", http.StatusUnauthorized}
	ErrMissingAuthData       = &APIError{"MISSING_AUTH_DATA", "Missing authentication headers", http.StatusUnauthorized}
	ErrClientNotFound        = &APIError{"CLIENT_NOT_FOUND", "API client not found", http.StatusUnauthorized}
	ErrClientInactive        = &APIError{"CLIENT_INACTIVE", "API client is inactive", http.StatusForbidden}
	ErrRateLimitExceeded     = &APIError{"RATE_LIMIT_EXCEEDED", "Too many requests, please try again later", http.StatusTooManyRequests}
	ErrInvalidRequest        = &APIError{"INVALID_REQUEST", "Invalid request format", http.StatusBadRequest}
	ErrValidationFailed      = &APIError{"VALIDATION_FAILED", "Request validation failed", http.StatusBadRequest}
	ErrInternalServerError   = &APIError{"INTERNAL_SERVER_ERROR", "An unexpected error occurred", http.StatusInternalServerError}
	ErrInvalidData           = &APIError{"INVALID_DATA", "Invalid data provided", http.StatusBadRequest}
	ErrDataTooLong           = &APIError{"DATA_TOO_LONG", "Data exceeds maximum length", http.StatusBadRequest}
	ErrRequiredField         = &APIError{"REQUIRED_FIELD_MISSING", "Required field is missing", http.StatusBadRequest}
	ErrRelatedRecordNotFound = &APIError{"RELATED_RECORD_NOT_FOUND", "Related record not found", http.StatusNotFound}
	ErrAlreadyExists         = &APIError{"ALREADY_EXISTS", "Resource already exists", http.StatusConflict}
	ErrRecordNotFound        = &APIError{"RECORD_NOT_FOUND", "Record not found", http.StatusNotFound}
	ErrWalletNotFound        = &APIError{"WALLET_NOT_FOUND", "Wallet not found", http.StatusNotFound}
	ErrEmptyAccountID        = &APIError{"EMPTY_ACCOUNT_ID", "Account ID cannot be empty", http.StatusBadRequest}
	ErrInvalidAccountID      = &APIError{"INVALID_ACCOUNT_ID", "Invalid account ID format", http.StatusBadRequest}
	ErrInsufficientFunds     = &APIError{"INSUFFICIENT_FUNDS", "Insufficient funds in wallet", http.StatusBadRequest}
	ErrInvalidWalletType     = &APIError{"INVALID_WALLET_TYPE", "Invalid wallet type", http.StatusBadRequest}
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
