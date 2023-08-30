package service

import (
	"context"
	"time"

	"github.com/PoorMercymain/user-segmenter/internal/domain"
)

type user struct {
	repo domain.UserRepository
}

func NewUser(repo domain.UserRepository) *user {
	return &user{repo: repo}
}

func (s *user) UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error {
	return s.repo.UpdateUserSegments(ctx, userID, slugsToAdd, slugsToDelete)
}

func (s *user) ReadUserSegments(ctx context.Context, userID string) ([]string, error) {
	return s.repo.ReadUserSegments(ctx, userID)
}

func (s *user) CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error {
	return s.repo.CreateDeletionTime(ctx, userID, slug, deletionTime)
}
