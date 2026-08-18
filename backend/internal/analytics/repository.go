package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	LogClick(ctx context.Context, shortCode, referrer, userAgent string) error
	GetTotalClicks(ctx context.Context, shortCode string) (int64, error)
	GetClickTimeline(ctx context.Context, shortCode string) ([]TimelinePoint, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) LogClick(ctx context.Context, shortCode, referrer, userAgent string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO clicks (url_id, referrer, user_agent)
	SELECT id, $2, $3 FROM urls WHERE short_code = $1`,
		shortCode, referrer, userAgent,
	)

	if err != nil {
		return fmt.Errorf("logging click: %w", err)
	}

	return nil
}

func(r *PostgresRepository) GetTotalClicks(ctx context.Context, shortCode string) (int64, error) {
	var count int64

	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM clicks c
		JOIN urls u ON u.id = c.url_id
		WHERE u.short_code = $1`,
		shortCode,
	).Scan(&count)

	if err != nil {
		return 0,fmt.Errorf("counting clicks: %w", err)
	}

	return count, nil
}


func (r *PostgresRepository) GetClickTimeline(ctx context.Context, shortCode string) ([]TimelinePoint, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DATE(c.clicked_at) AS day, COUNT(*) AS count
		 FROM clicks c
		 JOIN urls u ON u.id = c.url_id
		 WHERE u.short_code = $1
		 GROUP BY day
		 ORDER BY day`,
		shortCode,
	)
	if err != nil {
		return nil, fmt.Errorf("querying click timeline: %w", err)
	}
	defer rows.Close()

	var timeline []TimelinePoint
	for rows.Next() {
		var day time.Time
		var count int64
		if err := rows.Scan(&day, &count); err != nil {
			return nil, fmt.Errorf("scanning timeline row: %w", err)
		}
		timeline = append(timeline, TimelinePoint{
			Date:  day.Format("2006-01-02"),
			Count: count,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating timeline rows: %w", err)
	}

	return timeline, nil
}
