package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClickRepository struct {
	db *pgxpool.Pool
}

func NewClickRepository(db *pgxpool.Pool) *ClickRepository {
	return &ClickRepository{db: db}
}

func (r *ClickRepository) Insert(ctx context.Context, event *model.ClickEvent) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO click_events (url_id, clicked_at, ip_address, user_agent, referer, country)
		 VALUES ($1, $2, $3::inet, $4, $5, $6)`,
		event.URLID, event.ClickedAt, nullIfEmpty(event.IPAddress),
		event.UserAgent, event.Referer, nullIfEmpty(event.Country),
	)
	if err != nil {
		return fmt.Errorf("insert click event: %w", err)
	}
	return nil
}

func (r *ClickRepository) GetClicksOverTime(ctx context.Context, urlID int64, since time.Time) ([]model.DailyClickCount, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DATE(clicked_at) as date, COUNT(*) as count
		 FROM click_events
		 WHERE url_id = $1 AND clicked_at >= $2
		 GROUP BY DATE(clicked_at)
		 ORDER BY date ASC`,
		urlID, since,
	)
	if err != nil {
		return nil, fmt.Errorf("get clicks over time: %w", err)
	}
	defer rows.Close()

	var clicks []model.DailyClickCount
	for rows.Next() {
		var c model.DailyClickCount
		var date time.Time
		if err := rows.Scan(&date, &c.Count); err != nil {
			return nil, fmt.Errorf("scan click: %w", err)
		}
		c.Date = date.Format("2006-01-02")
		clicks = append(clicks, c)
	}
	return clicks, nil
}

func (r *ClickRepository) GetBrowserBreakdown(ctx context.Context, urlID int64, since time.Time) ([]model.BrowserCount, error) {
	rows, err := r.db.Query(ctx,
		`SELECT user_agent, COUNT(*) as count
		 FROM click_events
		 WHERE url_id = $1 AND clicked_at >= $2 AND user_agent IS NOT NULL
		 GROUP BY user_agent
		 ORDER BY count DESC`,
		urlID, since,
	)
	if err != nil {
		return nil, fmt.Errorf("get browser breakdown: %w", err)
	}
	defer rows.Close()

	// Raw user agents are returned; the service layer parses them into browser names
	var results []model.BrowserCount
	for rows.Next() {
		var b model.BrowserCount
		if err := rows.Scan(&b.Name, &b.Count); err != nil {
			return nil, fmt.Errorf("scan browser: %w", err)
		}
		results = append(results, b)
	}
	return results, nil
}

func (r *ClickRepository) GetCountryBreakdown(ctx context.Context, urlID int64, since time.Time) ([]model.CountryCount, error) {
	rows, err := r.db.Query(ctx,
		`SELECT COALESCE(country, 'Unknown') as country, COUNT(*) as count
		 FROM click_events
		 WHERE url_id = $1 AND clicked_at >= $2
		 GROUP BY country
		 ORDER BY count DESC`,
		urlID, since,
	)
	if err != nil {
		return nil, fmt.Errorf("get country breakdown: %w", err)
	}
	defer rows.Close()

	var results []model.CountryCount
	for rows.Next() {
		var c model.CountryCount
		if err := rows.Scan(&c.Code, &c.Count); err != nil {
			return nil, fmt.Errorf("scan country: %w", err)
		}
		results = append(results, c)
	}
	return results, nil
}

func (r *ClickRepository) GetTotalClicks(ctx context.Context, urlID int64, since time.Time) (int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events WHERE url_id = $1 AND clicked_at >= $2`,
		urlID, since,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get total clicks: %w", err)
	}
	return total, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
