package mapper

import (
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/infrastructure/database/models"
)

type ClientMapper struct{}

func NewClientMapper() *ClientMapper {
	return &ClientMapper{}
}

func (m *ClientMapper) ToDomain(dbClient *models.APIClient) *entity.APIClient {
	return &entity.APIClient{
		ID:        dbClient.ID,
		UserID:    dbClient.UserID,
		SecretKey: dbClient.SecretKey,
		IsActive:  dbClient.IsActive,
		CreatedAt: dbClient.CreatedAt,
		UpdatedAt: dbClient.UpdatedAt,
	}
}

func (m *ClientMapper) ToModel(client *entity.APIClient) *models.APIClient {
	return &models.APIClient{
		ID:        client.ID,
		UserID:    client.UserID,
		SecretKey: client.SecretKey,
		IsActive:  client.IsActive,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}
