package shortener

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	GetShortCodeURL(ctx context.Context, shortCode string) (URL, error)
	CreateShortCodeURL(ctx context.Context, OriginalURL string) (URL, error)
	UpdateShortCodeURL(ctx context.Context, id int64, shortCode string) error
}
 
// handles queries
// 1. -- getting the shortcode of a url 2. entering/inserting a row of a new url 3. updating the row once the new code is computed


func NewPosgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) GetShortCodeURL(ctx context.Context, shortCode string) (*URL, error) {
	var url URL
	err := r.pool.QueryRow(ctx, 
		`SELECT * FROM urls WHERE short_code = $1`, 
		shortCode,
	).Scan(&url.Id, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &url, err
} 