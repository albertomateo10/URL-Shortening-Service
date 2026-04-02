package service

import (
	"context"
	"fmt"
	"time"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/albertomateo10/url-shortener/backend/internal/repository"
	"github.com/mssola/useragent"
)

type AnalyticsService struct {
	clickRepo *repository.ClickRepository
	urlRepo   *repository.URLRepository
}

func NewAnalyticsService(clickRepo *repository.ClickRepository, urlRepo *repository.URLRepository) *AnalyticsService {
	return &AnalyticsService{clickRepo: clickRepo, urlRepo: urlRepo}
}

func (s *AnalyticsService) GetClicksOverTime(ctx context.Context, urlID int64, period string) (*model.ClicksOverTimeResponse, error) {
	u, err := s.urlRepo.GetByID(ctx, urlID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("url not found")
	}

	since, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}

	clicks, err := s.clickRepo.GetClicksOverTime(ctx, urlID, since)
	if err != nil {
		return nil, err
	}
	if clicks == nil {
		clicks = []model.DailyClickCount{}
	}

	total, err := s.clickRepo.GetTotalClicks(ctx, urlID, since)
	if err != nil {
		return nil, err
	}

	return &model.ClicksOverTimeResponse{
		URLID:        urlID,
		Period:       period,
		TotalClicks:  total,
		ClicksPerDay: clicks,
	}, nil
}

func (s *AnalyticsService) GetSources(ctx context.Context, urlID int64, period string) (*model.SourcesResponse, error) {
	u, err := s.urlRepo.GetByID(ctx, urlID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("url not found")
	}

	since, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}

	rawBrowsers, err := s.clickRepo.GetBrowserBreakdown(ctx, urlID, since)
	if err != nil {
		return nil, err
	}

	// Parse raw user-agent strings into browser names and aggregate
	browsers := aggregateBrowsers(rawBrowsers)

	countries, err := s.clickRepo.GetCountryBreakdown(ctx, urlID, since)
	if err != nil {
		return nil, err
	}
	if countries == nil {
		countries = []model.CountryCount{}
	}

	return &model.SourcesResponse{
		URLID:     urlID,
		Period:    period,
		Browsers:  browsers,
		Countries: countries,
	}, nil
}

func aggregateBrowsers(raw []model.BrowserCount) []model.BrowserCount {
	counts := make(map[string]int)
	for _, b := range raw {
		ua := useragent.New(b.Name)
		name, _ := ua.Browser()
		if name == "" {
			name = "Other"
		}
		counts[name] += b.Count
	}

	result := make([]model.BrowserCount, 0, len(counts))
	for name, count := range counts {
		result = append(result, model.BrowserCount{Name: name, Count: count})
	}
	return result
}

func parsePeriod(period string) (time.Time, error) {
	now := time.Now()
	switch period {
	case "24h":
		return now.Add(-24 * time.Hour), nil
	case "7d":
		return now.AddDate(0, 0, -7), nil
	case "30d":
		return now.AddDate(0, 0, -30), nil
	case "90d":
		return now.AddDate(0, 0, -90), nil
	default:
		return time.Time{}, fmt.Errorf("invalid period: %s (use 24h, 7d, 30d, or 90d)", period)
	}
}
