package database

import (
	"e-wallet/internal/infrastructure/database/models"
	"e-wallet/internal/infrastructure/logger"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	logger.Info.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.APIClient{},
		&models.Wallet{},
		&models.Transaction{},
	)
	if err != nil {
		return err
	}

	logger.Info.Println("Database migrations completed successfully")
	return nil
}
