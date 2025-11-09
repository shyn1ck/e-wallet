package middleware

import (
	"e-wallet/internal/delivery/http/handler"
	"e-wallet/internal/infrastructure/logger"
	apperrors "e-wallet/pkg/errors"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")
				logger.Error.Printf("[%s] Panic recovered: %v", requestID, err)

				handler.HandleError(c, apperrors.ErrInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
