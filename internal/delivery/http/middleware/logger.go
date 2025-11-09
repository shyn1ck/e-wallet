package middleware

import (
	"e-wallet/internal/infrastructure/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		requestID := c.GetString("request_id")

		logger.Info.Printf("[%s] %s %s %d %v | IP: %s",
			requestID, method, path, statusCode, duration, clientIP)
	}
}
