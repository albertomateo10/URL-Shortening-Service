package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/redis/go-redis/v9"
)

const cacheTTL = 24 * time.Hour

type CacheService struct {
	rdb *redis.Client
}

func NewCacheService(rdb *redis.Client) *CacheService {
	return &CacheService{rdb: rdb}
}

func (s *CacheService) GetURL(ctx context.Context, shortCode string) (*model.URL, error) {
	val, err := s.rdb.Get(ctx, "url:"+shortCode).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var u model.URL
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *CacheService) SetURL(ctx context.Context, u *model.URL) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, "url:"+u.ShortCode, data, cacheTTL).Err()
}

func (s *CacheService) DeleteURL(ctx context.Context, shortCode string) error {
	return s.rdb.Del(ctx, "url:"+shortCode).Err()
}
