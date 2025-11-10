package usecase

import (
	"context"
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/repository"
	"e-wallet/internal/infrastructure/cache"
	"e-wallet/internal/infrastructure/logger"
	"encoding/json"
)

type ClientCacheUseCase struct {
	clientRepo repository.ClientRepository
	cacheRepo  repository.CacheRepository
}

func NewClientCacheUseCase(
	clientRepo repository.ClientRepository,
	cacheRepo repository.CacheRepository,
) *ClientCacheUseCase {
	return &ClientCacheUseCase{
		clientRepo: clientRepo,
		cacheRepo:  cacheRepo,
	}
}

func (uc *ClientCacheUseCase) GetClient(ctx context.Context, userID string) (*entity.APIClient, error) {
	cacheKey := cache.BuildAPIClientKey(userID)

	cachedData, err := uc.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var client entity.APIClient
		if err := json.Unmarshal([]byte(cachedData), &client); err == nil {
			logger.Info.Printf("[ClientCacheUseCase.GetClient]: Cache HIT for user_id: %s", userID)
			return &client, nil
		}
		logger.Info.Printf("[ClientCacheUseCase.GetClient]: Failed to deserialize cached client for user_id: %s", userID)
	}

	logger.Info.Printf("[ClientCacheUseCase.GetClient]: Cache MISS for user_id: %s, querying database", userID)
	client, err := uc.clientRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	clientJSON, err := json.Marshal(client)
	if err == nil {
		if err := uc.cacheRepo.Set(ctx, cacheKey, string(clientJSON), cache.TTLAPIClient); err != nil {
			logger.Error.Printf("[ClientCacheUseCase.GetClient]: Failed to cache client for user_id: %s: %v", userID, err)
		}
	}

	return client, nil
}

func (uc *ClientCacheUseCase) InvalidateClient(ctx context.Context, userID string) error {
	cacheKey := cache.BuildAPIClientKey(userID)
	if err := uc.cacheRepo.Delete(ctx, cacheKey); err != nil {
		logger.Error.Printf("[ClientCacheUseCase.InvalidateClient]: Failed to invalidate cache for user_id: %s: %v", userID, err)
		return err
	}
	logger.Info.Printf("[ClientCacheUseCase.InvalidateClient]: Cache invalidated for user_id: %s", userID)
	return nil
}
