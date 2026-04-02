package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const cacheTTL = 24 * time.Hour

type CacheService struct {
	rdb *redis.Client
}

func NewCacheService(rdb *redis.Client) *CacheService {
	return &CacheService{rdb: rdb}
}

func (s *CacheService) GetURL(ctx context.Context, shortCode string) (string, error) {
	val, err := s.rdb.Get(ctx, "url:"+shortCode).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *CacheService) SetURL(ctx context.Context, shortCode, originalURL string) error {
	return s.rdb.Set(ctx, "url:"+shortCode, originalURL, cacheTTL).Err()
}

func (s *CacheService) DeleteURL(ctx context.Context, shortCode string) error {
	return s.rdb.Del(ctx, "url:"+shortCode).Err()
}
