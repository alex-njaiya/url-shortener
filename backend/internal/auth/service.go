package auth

import (
	"context"
	"errors"
	"fmt"
)


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
	
	user, err := s.repo.GetUserByEmail(ctx, email)

	if errors.Is(err, ErrNotFound) {
		// meaning the user does not exist
		// hash the password then store
		hashedPassword, err := HashPassword(password)

		if err != nil {
			return nil, err
		}

		user, err = s.repo.InsertUser(ctx, firstName, lastName, email, hashedPassword)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	// user already exists
	fmt.Println("User already exists. Please login")

	return user, nil
}