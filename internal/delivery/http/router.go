package http

import (
	"e-wallet/internal/delivery/http/handler"
	"e-wallet/internal/delivery/http/middleware"
	"e-wallet/internal/domain/repository"
	"e-wallet/pkg/crypto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	WalletHandler *handler.WalletHandler
	ClientRepo    repository.ClientRepository
	HMACAlgorithm crypto.HMACAlgorithm
	GinMode       string
}

func NewRouter(cfg *RouterConfig) *gin.Engine {
	// Set Gin mode based on configuration
	switch cfg.GinMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.DebugMode) // default to debug
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())

	// Rate limiter: 100 requests per minute
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	router.Use(rateLimiter.Middleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "e-wallet-api",
		})
	})

	// API v1 routes with HMAC authentication
	v1 := router.Group("/api/v1")
	v1.Use(middleware.HMACAuth(cfg.ClientRepo, cfg.HMACAlgorithm))
	{
		// Wallet routes
		wallet := v1.Group("/wallet")
		{
			wallet.POST("/check", cfg.WalletHandler.CheckWallet)
			wallet.POST("/deposit", cfg.WalletHandler.Deposit)
			wallet.POST("/balance", cfg.WalletHandler.GetBalance)
			wallet.POST("/monthly-stats", cfg.WalletHandler.GetMonthlyStats)
		}
	}

	return router
}
