package repository

import (
	"context"
	"fmt"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(ctx context.Context, shortCode, originalURL string) (*model.URL, error) {
	var u model.URL
	err := r.db.QueryRow(ctx,
		`INSERT INTO urls (short_code, original_url) VALUES ($1, $2)
		 RETURNING id, short_code, original_url, click_count, created_at`,
		shortCode, originalURL,
	).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create url: %w", err)
	}
	return &u, nil
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	var u model.URL
	err := r.db.QueryRow(ctx,
		`SELECT id, short_code, original_url, click_count, created_at
		 FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get url by short code: %w", err)
	}
	return &u, nil
}

func (r *URLRepository) GetByID(ctx context.Context, id int64) (*model.URL, error) {
	var u model.URL
	err := r.db.QueryRow(ctx,
		`SELECT id, short_code, original_url, click_count, created_at
		 FROM urls WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get url by id: %w", err)
	}
	return &u, nil
}

func (r *URLRepository) List(ctx context.Context, page, limit int) ([]model.URL, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM urls`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count urls: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, short_code, original_url, click_count, created_at
		 FROM urls ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list urls: %w", err)
	}
	defer rows.Close()

	var urls []model.URL
	for rows.Next() {
		var u model.URL
		if err := rows.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan url: %w", err)
		}
		urls = append(urls, u)
	}
	return urls, total, nil
}

func (r *URLRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM urls WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("url not found")
	}
	return nil
}

func (r *URLRepository) IncrementClickCount(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE urls SET click_count = click_count + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("increment click count: %w", err)
	}
	return nil
}
