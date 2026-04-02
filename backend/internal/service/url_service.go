package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/albertomateo10/url-shortener/backend/internal/repository"
	"github.com/albertomateo10/url-shortener/backend/internal/shortcode"
)

const maxCollisionRetries = 3

type URLService struct {
	repo    *repository.URLRepository
	cache   *CacheService
	baseURL string
}

func NewURLService(repo *repository.URLRepository, cache *CacheService, baseURL string) *URLService {
	return &URLService{repo: repo, cache: cache, baseURL: baseURL}
}

func (s *URLService) CreateURL(ctx context.Context, rawURL string) (*model.CreateURLResponse, error) {
	if err := validateURL(rawURL); err != nil {
		return nil, err
	}

	var u *model.URL
	for i := 0; i < maxCollisionRetries; i++ {
		code, err := shortcode.Generate()
		if err != nil {
			return nil, fmt.Errorf("generate short code: %w", err)
		}

		u, err = s.repo.Create(ctx, code, rawURL)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			return nil, err
		}
		break
	}
	if u == nil {
		return nil, fmt.Errorf("failed to generate unique short code after %d attempts", maxCollisionRetries)
	}

	// Pre-populate cache
	_ = s.cache.SetURL(ctx, u.ShortCode, u.OriginalURL)

	return &model.CreateURLResponse{
		ID:          u.ID,
		ShortCode:   u.ShortCode,
		OriginalURL: u.OriginalURL,
		ShortURL:    fmt.Sprintf("%s/r/%s", s.baseURL, u.ShortCode),
		CreatedAt:   u.CreatedAt,
	}, nil
}

func (s *URLService) GetURL(ctx context.Context, id int64) (*model.URLResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return s.toResponse(u), nil
}

func (s *URLService) ListURLs(ctx context.Context, page, limit int) (*model.URLListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	urls, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]model.URLResponse, len(urls))
	for i, u := range urls {
		responses[i] = *s.toResponse(&u)
	}

	return &model.URLListResponse{
		URLs:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *URLService) DeleteURL(ctx context.Context, id int64) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return fmt.Errorf("url not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.DeleteURL(ctx, u.ShortCode)
	return nil
}

func (s *URLService) ResolveShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	// // Try cache first
	// cached, err := s.cache.GetURL(ctx, shortCode)
	// if err == nil && cached != "" {
	// 	return &model.URL{ShortCode: shortCode, OriginalURL: cached}, nil
	// }

	// Cache miss — query DB
	u, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}

	// Populate cache
	_ = s.cache.SetURL(ctx, u.ShortCode, u.OriginalURL)

	return u, nil
}

func (s *URLService) toResponse(u *model.URL) *model.URLResponse {
	return &model.URLResponse{
		ID:          u.ID,
		ShortCode:   u.ShortCode,
		OriginalURL: u.OriginalURL,
		ShortURL:    fmt.Sprintf("%s/r/%s", s.baseURL, u.ShortCode),
		ClickCount:  u.ClickCount,
		CreatedAt:   u.CreatedAt,
	}
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("url is required")
	}
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}
