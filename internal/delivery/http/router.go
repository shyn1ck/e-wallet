package http

import (
	"e-wallet/internal/delivery/http/handler"
	"e-wallet/internal/delivery/http/middleware"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/infrastructure/config"
	"e-wallet/internal/usecase"
	"e-wallet/pkg/crypto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "e-wallet/docs"
)

type RouterConfig struct {
	WalletHandler      *handler.WalletHandler
	ClientRepo         repository.ClientRepository
	ClientCacheUseCase *usecase.ClientCacheUseCase
	HMACAlgorithm      crypto.HMACAlgorithm
	GinMode            string
	Environment        string
}

func NewRouter(cfg *RouterConfig) *gin.Engine {
	switch cfg.GinMode {
	case config.GinModeRelease:
		gin.SetMode(gin.ReleaseMode)
	case config.GinModeDebug:
		gin.SetMode(gin.DebugMode)
	case config.GinModeTest:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
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

	// Swagger documentation (only in development)
	if cfg.Environment == config.EnvironmentDevelopment {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 routes with HMAC authentication
	v1 := router.Group("/api/v1")
	v1.Use(middleware.HMACAuth(cfg.ClientRepo, cfg.ClientCacheUseCase, cfg.HMACAlgorithm))
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
