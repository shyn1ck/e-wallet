package container

import (
	"e-wallet/internal/delivery/http"
	"e-wallet/internal/delivery/http/handler"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/domain/service"
	"e-wallet/internal/infrastructure/cache"
	"e-wallet/internal/infrastructure/config"
	"e-wallet/internal/infrastructure/database"
	"e-wallet/internal/repository/postgres"
	"e-wallet/internal/usecase"
	"e-wallet/pkg/crypto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Container struct {
	Config *config.Config
	DB     *gorm.DB
	Cache  *cache.RedisClient

	// Repositories
	WalletRepo      repository.WalletRepository
	TransactionRepo repository.TransactionRepository
	ClientRepo      repository.ClientRepository

	// Services
	BalanceValidator *service.BalanceValidator

	// Use Cases
	WalletCheckUseCase        *usecase.WalletCheckUseCase
	WalletDepositUseCase      *usecase.WalletDepositUseCase
	WalletBalanceUseCase      *usecase.WalletBalanceUseCase
	WalletMonthlyStatsUseCase *usecase.WalletMonthlyStatsUseCase

	// Handlers
	WalletHandler *handler.WalletHandler

	// Router
	Router *gin.Engine
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config: cfg,
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		return nil, err
	}
	c.DB = db

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		return nil, err
	}

	// Initialize cache
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, err
	}
	c.Cache = redisClient

	// Initialize repositories
	c.WalletRepo = postgres.NewWalletRepository(db)
	c.TransactionRepo = postgres.NewTransactionRepository(db)
	c.ClientRepo = postgres.NewClientRepository(db)

	// Initialize domain services
	c.BalanceValidator = service.NewBalanceValidator()

	// Initialize use cases
	c.WalletCheckUseCase = usecase.NewWalletCheckUseCase(c.WalletRepo)
	c.WalletDepositUseCase = usecase.NewWalletDepositUseCase(
		db,
		c.WalletRepo,
		c.TransactionRepo,
		c.BalanceValidator,
	)
	c.WalletBalanceUseCase = usecase.NewWalletBalanceUseCase(c.WalletRepo)
	c.WalletMonthlyStatsUseCase = usecase.NewWalletMonthlyStatsUseCase(
		c.WalletRepo,
		c.TransactionRepo,
	)

	// Initialize handlers
	c.WalletHandler = handler.NewWalletHandler(
		c.WalletCheckUseCase,
		c.WalletDepositUseCase,
		c.WalletBalanceUseCase,
		c.WalletMonthlyStatsUseCase,
	)

	// Initialize router
	c.Router = http.NewRouter(&http.RouterConfig{
		WalletHandler: c.WalletHandler,
		ClientRepo:    c.ClientRepo,
		HMACAlgorithm: crypto.HMACAlgorithm(cfg.Auth.HMACAlgorithm),
		GinMode:       cfg.App.GinMode,
	})

	return c, nil
}

func (c *Container) Close() error {
	if c.Cache != nil {
		if err := c.Cache.Close(); err != nil {
			return err
		}
	}

	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return nil
}
