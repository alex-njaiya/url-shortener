package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	GetShortCodeURL(ctx context.Context, shortCode string) (*URL, error)
	InsertURL(ctx context.Context, userId *int64, shortCode, originalURL string) (*URL, error)
	GetURLsByUserID(ctx context.Context, userId int64) ([]*URL, error)
}

var ErrNotFound = errors.New("short url record not found")

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) GetShortCodeURL(ctx context.Context, shortCode string) (*URL, error) {
	url := new(URL)
	err := r.pool.QueryRow(ctx,
		`SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&url.Id, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fetching url by code: %w", err)
	}

	return url, err
}

func (r *PostgresRepository) InsertURL(ctx context.Context, userId *int64, shortCode, originalURL string) (*URL, error) {
	url := new(URL)

	err := r.pool.QueryRow(ctx,
		`INSERT INTO urls (user_id, short_code, original_url)
		 VALUES ($1, $2, $3) RETURNING id, short_code, original_url, created_at`,
		userId, shortCode, originalURL,
	).Scan(&url.Id, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("adding new url: %w", err)
	}

	return url, nil
}

func (r *PostgresRepository) GetURLsByUserID(ctx context.Context, userID int64) ([]*URL, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, short_code, original_url, created_at FROM urls
		 WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("fetching urls for user: %w", err)
	}
	defer rows.Close()

	var urls []*URL
	for rows.Next() {
		u := new(URL)
		if err := rows.Scan(&u.Id, &u.ShortCode, &u.OriginalURL, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning url row: %w", err)
		}
		urls = append(urls, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating url rows: %w", err)
	}

	return urls, nil
}
