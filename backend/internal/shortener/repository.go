package shortener

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	GetShortCodeURL(ctx context.Context, shortCode string) (*URL, error)
	InsertURL(ctx context.Context, id int64, shortCode string, OriginalURL string) (time.Time, error)
	ReserveID(ctx context.Context) (int64, error)
}

// handles queries
// 1. -- getting the shortcode of a url 2. entering/inserting a row of a new url 3. updating the row once the new code is computed

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) GetShortCodeURL(ctx context.Context, shortCode string) (*URL, error) {
	url := new(URL)
	err := r.pool.QueryRow(ctx,
		`SELECT * FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&url.Id, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("fetching url by code: %w", err)
	}
	return url, err
}

func (r *PostgresRepository) InsertURL(ctx context.Context, id int64, shortCode string, originalURL string) (time.Time, error) {
	// instead of using the conventional design of transaction rollbacks we will use the sequence to the get the int
	// for the next row in the table
	var timestamp time.Time

	err := r.pool.QueryRow(ctx,
		`INSERT INTO urls (id, short_code, original_url) VALUES ($1, $2, $3)
		RETURNING created_at`,
		id, shortCode, originalURL,
	).Scan(&timestamp)

	if err != nil {
		return time.Time{}, fmt.Errorf("inserting url: %w", err)
	}

	return timestamp, nil
}

func (r *PostgresRepository) ReserveID(ctx context.Context) (int64, error) {
	var sequenceName string

	err := r.pool.QueryRow(ctx,
		`SELECT pg_get_serial_sequence($1, $2)`,
		"urls", "id",
	).Scan(&sequenceName)

	if err != nil {
		return 0, fmt.Errorf("looking up sequence name: %w", err)
	}

	var nextID int64

	err = r.pool.QueryRow(ctx,
		`SELECT nextval($1)`,
		sequenceName,
	).Scan(&nextID)

	if err != nil {
		return 0, fmt.Errorf("reserving next id: %w", err)
	}

	return nextID, nil
}
