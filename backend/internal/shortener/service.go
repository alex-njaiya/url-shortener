package shortener

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const predefinedPreffix = "alx"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ShortenUsingBase62Encode(ctx context.Context, originalURL string) (*URL, error) {
	url := new(URL)
	// check if the original url has a prefix of either http or https
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		return nil, fmt.Errorf("original url must start with http:// or https://")
	}

	id, err := s.repo.ReserveID(ctx)

	if err != nil {
		return nil, fmt.Errorf("getting the nextval Id: %w", err)
	}

	// generate shortcode
	code := base62encode(id)

	timestamp, err := s.repo.InsertURL(ctx, id, code, originalURL)

	if err != nil {
		return nil, fmt.Errorf("creating short url: %w", err)
	}

	url.Id = id
	url.ShortCode = code
	url.OriginalURL = originalURL
	url.CreatedAt = timestamp

	return url, nil
}

func (s *Service) ShortenByHashing(ctx context.Context, userId *int64, originalURL string) (*URL, error) {
	hashCode := hashLongUrl(originalURL)

	// check whether the hash code exists on db
	url, err := s.repo.GetShortCodeURL(ctx, hashCode)
	
	// // check for system/db errors
	// if err != nil {
	// 	return nil, fmt.Errorf("service failed: %w", err)
	// }

	// check if it does not exist
	if !errors.Is(err, ErrNotFound) {
		// logic for when the url with the generated code does exist. Meaning there is a duplicate
		// there is a collision we need to define a collision strategy
		newCode := predefinedPreffix + hashCode

		// with a new hash code insert into the db
		url, err := s.repo.InsertURLWhenUsingHashing(ctx, userId, newCode, originalURL)

		if err != nil {
			return nil, err
		}

		return url, nil

	}


	// if the short code does not exist in the db then insert as is in the db
	url, err = s.repo.InsertURLWhenUsingHashing(ctx, userId, hashCode, originalURL)

	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *Service) Resolve(ctx context.Context, code string) (*URL, error) {
	url, err := s.repo.GetShortCodeURL(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("getting original url from shortcode: %w", err)
	}

	return url, nil
}

func base62encode(n int64) string {
	if n == 0 {
		return string(base62Alphabet[0])
	}

	var sb strings.Builder
	base := int64(len(base62Alphabet))

	for n > 0 {
		sb.WriteByte(base62Alphabet[n%base])
		n /= base
	}

	s := sb.String()
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func hashLongUrl(originalURL string) string {
	hashByte := sha256.Sum256([]byte(originalURL))

	hashStr := hex.EncodeToString(hashByte[:])

	return hashStr[:6]
}
