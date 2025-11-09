package errors

import (
	stderr "errors"

	"e-wallet/internal/infrastructure/logger"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// TranslateError translates PostgreSQL and GORM errors to domain errors
func TranslateError(err error) error {
	if err == nil {
		return nil
	}

	// Check for GORM record not found
	if stderr.Is(err, gorm.ErrRecordNotFound) {
		logger.Info.Printf("Record not found error: %v", err)
		return ErrRecordNotFound
	}

	// Check for PostgreSQL errors
	var pgErr *pgconn.PgError
	if stderr.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			logger.Error.Printf("Uniqueness violation: %v", err)
			return ErrAlreadyExists

		case "23503": // foreign_key_violation
			logger.Error.Printf("Foreign key violation: %v", err)
			return ErrRelatedRecordNotFound

		case "23502": // not_null_violation
			logger.Error.Printf("Not null constraint violation: %v", err)
			return ErrRequiredField

		case "22001": // string_data_right_truncation
			logger.Error.Printf("String too long error: %v", err)
			return ErrDataTooLong

		case "23514": // check_violation
			logger.Error.Printf("Check constraint violation: %v", err)
			return ErrInvalidData

		case "40P01": // deadlock_detected
			logger.Error.Printf("Deadlock detected: %v", err)
			return ErrInternalServerError

		case "42702": // ambiguous_column
			logger.Error.Printf("Ambiguous column error: %v", err)
			return ErrInternalServerError

		default:
			logger.Error.Printf("Unhandled PostgreSQL error (code %s): %v", pgErr.Code, err)
			return ErrInternalServerError
		}
	}

	logger.Warning.Printf("Unhandled error: %v", err)
	return ErrInternalServerError
}
