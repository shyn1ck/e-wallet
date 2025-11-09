package middleware

import (
	"e-wallet/internal/delivery/http/handler"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/infrastructure/logger"
	"e-wallet/pkg/crypto"
	apperrors "e-wallet/pkg/errors"
	"io"

	"github.com/gin-gonic/gin"
)

// HMACAuth validates HMAC authentication
func HMACAuth(clientRepo repository.ClientRepository, algorithm crypto.HMACAlgorithm) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-UserId")
		digest := c.GetHeader("X-Digest")

		if userID == "" || digest == "" {
			logger.Warning.Printf("[middleware.HMACAuth]: Missing authentication headers")
			handler.HandleError(c, apperrors.ErrMissingAuthData)
			c.Abort()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error.Printf("[middleware.HMACAuth]: Failed to read request body: %v", err)
			handler.HandleError(c, apperrors.ErrInvalidRequest)
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(&bodyReader{body: body, pos: 0})

		client, err := clientRepo.FindByUserID(c.Request.Context(), userID)
		if err != nil {
			logger.Error.Printf("[middleware.HMACAuth]: Client not found: %s", userID)
			handler.HandleError(c, apperrors.ErrClientNotFound)
			c.Abort()
			return
		}

		if !client.IsActive {
			logger.Warning.Printf("[middleware.HMACAuth]: Inactive client attempted access: %s", userID)
			handler.HandleError(c, apperrors.ErrClientInactive)
			c.Abort()
			return
		}

		if !crypto.ValidateHMAC(algorithm, client.SecretKey, string(body), digest) {
			logger.Warning.Printf("[middleware.HMACAuth]: Invalid HMAC signature for user: %s", userID)
			handler.HandleError(c, apperrors.ErrInvalidSignature)
			c.Abort()
			return
		}

		c.Set("client_id", client.ID)
		c.Set("user_id", userID)

		c.Next()
	}
}

type bodyReader struct {
	body []byte
	pos  int
}

func (br *bodyReader) Read(p []byte) (n int, err error) {
	if br.pos >= len(br.body) {
		return 0, io.EOF
	}
	n = copy(p, br.body[br.pos:])
	br.pos += n
	return n, nil
}
