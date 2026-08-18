package analytics

import (
	"context"
	"fmt"
	"log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}


func (s *Service) LogClick(code, referrer, userAgent string) {
	ctx := context.Background()
	if err := s.repo.LogClick(ctx, code, referrer, userAgent); err != nil {
		log.Printf("analytics: failed to log click for %s: %v", code, err)
	}
}

func (s *Service) Stats(ctx context.Context, shortCode string) (*Stats, error) {
	total, err := s.repo.GetTotalClicks(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("getting total clicks: %w", err)
	}

	timeline, err := s.repo.GetClickTimeline(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("getting click timeline: %w", err)
	}

	return &Stats{
		ShortCode:   shortCode,
		TotalClicks: total,
		Timeline:    timeline,
	}, nil
}

