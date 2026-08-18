package auth

import (
	"context"
	"errors"
	"fmt"
)

var ErrUserExists = errors.New("a user with this email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterUser(ctx context.Context, firstName, lastName, email, password string) (*User, error) {
	// check if the user exists first
	// If the user exists return an error that the user already exists
	// If the user does not exist insert the user into the database and give them access
	// Before inserting them into the db hash the password first

	_, err := s.repo.GetUserByEmail(ctx, email)

	if err == nil {
		return nil, ErrUserExists
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}

	hashedPassword, err := HashPassword(password)

	if err != nil {
		return nil, err
	}

	user, err := s.repo.InsertUser(ctx, firstName, lastName, email, hashedPassword)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if errors.Is(err, ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("looking up user: %w", err)
	}

	if err := CheckPassword(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
