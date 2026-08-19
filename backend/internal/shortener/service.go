package shortener

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

const predefinedPreffix = "alx"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ShortenByHashing(ctx context.Context, userId *int64, originalURL string) (*URL, error) {
	hashCode := hashLongUrl(originalURL)

	// check whether the hash code exists on db
	_, err := s.repo.GetShortCodeURL(ctx, hashCode)

	switch {
	case err == nil:
		// Collision: this hash is already in use by another URL.
		newCode := predefinedPreffix + hashCode
		url, err := s.repo.InsertURL(ctx, userId, newCode, originalURL)
		if err != nil {
			return nil, err
		}
		return url, nil

	case errors.Is(err, ErrNotFound):
		// No collision -- free to use.
		url, err := s.repo.InsertURL(ctx, userId, hashCode, originalURL)
		if err != nil {
			return nil, err
		}
		return url, nil

	default:
		return nil, fmt.Errorf("checking for existing short code: %w", err)
	}
}

func (s *Service) Resolve(ctx context.Context, code string) (*URL, error) {
	url, err := s.repo.GetShortCodeURL(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("getting original url from shortcode: %w", err)
	}

	return url, nil
}


func (s *Service) GetUserURLs(ctx context.Context, userID int64) ([]*URL, error) {
	urls, err := s.repo.GetURLsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting urls for user: %w", err)
	}
	return urls, nil
}

func hashLongUrl(originalURL string) string {
	hashByte := sha256.Sum256([]byte(originalURL))

	hashStr := hex.EncodeToString(hashByte[:])

	return hashStr[:6]
}
