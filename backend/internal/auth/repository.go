package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	InsertUser(ctx context.Context, firstName, lastName, email, password string) (*User, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

var ErrNotFound = errors.New("User not found")


func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}


func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := new(User)

	err := r.pool.QueryRow(ctx, 
		`SELECT id, first_name, last_name, email, password, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fetching user by email: %w",err)
	}

	return user, nil
}


func (r *PostgresRepository) InsertUser(ctx context.Context, firstname, lastname, email, hashedPassword string) (*User, error) {
	user := &User{
		FirstName: firstname,
		LastName: lastname,
		Email: email,
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)
		RETURNING id, first_name, last_name, email, password, created_at`,
		firstname, lastname, email, hashedPassword,
	).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("inserting user into db: %w", err)
	}

	return user, nil
}