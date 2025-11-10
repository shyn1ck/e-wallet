package handler

import (
	"e-wallet/internal/dto/response"
	apperrors "e-wallet/pkg/errors"

	"github.com/gin-gonic/gin"
)

// HandleError handles errors and returns appropriate HTTP response
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	statusCode := apperrors.GetStatusCode(err)
	errorCode := apperrors.GetErrorCode(err)
	errorMessage := apperrors.GetErrorMessage(err)

	c.JSON(statusCode, response.ErrorResponse{
		Error:   errorCode,
		Message: errorMessage,
	})
}
