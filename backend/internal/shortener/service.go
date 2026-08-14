package shortener

import (
	"context"
	"fmt"
	"strings"
)

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type Service struct {
	repo Repository
}


func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}


func (s *Service) Shorten(ctx context.Context, originalURL string) (*URL, error) {
	url := new(URL)
	// check if the original url has a prefix of either http or https
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		return nil, fmt.Errorf("original url must start with http or https")
	}

	id, err := s.repo.ReserveID(ctx) 

	if err != nil {
		return nil, fmt.Errorf("getting the nextval Id: %w", err)
	}

	//generate the shortcode
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

	for i, j := 0, len(runes) - 1; i < j; i , j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}