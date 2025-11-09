package repository

import (
	"context"
	"e-wallet/internal/domain/entity"
)

// ClientRepository defines the interface for client persistence
type ClientRepository interface {
	FindByUserID(ctx context.Context, userID string) (*entity.APIClient, error)
	Create(ctx context.Context, client *entity.APIClient) error
	Update(ctx context.Context, client *entity.APIClient) error
}
