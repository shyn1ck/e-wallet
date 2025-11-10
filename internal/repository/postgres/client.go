package postgres

import (
	"context"
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/infrastructure/database/models"
	"e-wallet/internal/infrastructure/logger"
	"e-wallet/internal/repository/mapper"
	apperrors "e-wallet/pkg/errors"
	"errors"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db     *gorm.DB
	mapper *mapper.ClientMapper
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db:     db,
		mapper: mapper.NewClientMapper(),
	}
}

// FindByUserID retrieves an API client by user ID
func (r *ClientRepository) FindByUserID(ctx context.Context, userID string) (*entity.APIClient, error) {
	var dbClient models.APIClient
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&dbClient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrClientNotFound
		}
		logger.Error.Printf("[postgres.FindByUserID]: Failed to find client by user_id '%s': %v", userID, err)
		return nil, apperrors.TranslateError(err)
	}

	return r.mapper.ToDomain(&dbClient), nil
}

// Create creates a new API client
func (r *ClientRepository) Create(ctx context.Context, client *entity.APIClient) error {
	dbClient := r.mapper.ToModel(client)
	err := r.db.WithContext(ctx).Create(dbClient).Error
	if err != nil {
		logger.Error.Printf("[postgres.Create]: Failed to create client: %v", err)
		return apperrors.TranslateError(err)
	}

	client.ID = dbClient.ID
	client.CreatedAt = dbClient.CreatedAt
	client.UpdatedAt = dbClient.UpdatedAt

	return nil
}

// Update updates an existing API client
func (r *ClientRepository) Update(ctx context.Context, client *entity.APIClient) error {
	dbClient := r.mapper.ToModel(client)
	err := r.db.WithContext(ctx).Save(dbClient).Error
	if err != nil {
		logger.Error.Printf("[postgres.Update]: Failed to update client id %d: %v", client.ID, err)
		return apperrors.TranslateError(err)
	}
	return nil
}
