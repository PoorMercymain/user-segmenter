package service

import (
	"context"
	"time"

	"github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	slugvalidator "github.com/PoorMercymain/user-segmenter/pkg/slug-validator"
)

type segment struct {
	repo domain.SegmentRepository
}

func NewSegment(repo domain.SegmentRepository) *segment {
	return &segment{repo: repo}
}

func (s *segment) CreateSegment(ctx context.Context, slug string) error {
	if slugvalidator.IsSlug(slug) {
		return s.repo.CreateSegment(ctx, slug)
	}

	return errors.ErrorNotASlug
}

func (s *segment) DeleteSegment(ctx context.Context, slug string) error {
	return s.repo.DeleteSegment(ctx, slug)
}

func (s *segment) UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error {
	return s.repo.UpdateUserSegments(ctx, userID, slugsToAdd, slugsToDelete)
}

func (s *segment) ReadUserSegments(ctx context.Context, userID string) ([]string, error) {
	return s.repo.ReadUserSegments(ctx, userID)
}

func (s *segment) ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time) ([]domain.HistoryElem, error) {
	return s.repo.ReadUserSegmentsHistory(ctx, userID, startDate, endDate)
}

func (s *segment) CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error {
	return s.repo.CreateDeletionTime(ctx, userID, slug, deletionTime)
}

func (s *segment) AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error {
	return s.repo.AddSegmentToPercentOfUsers(ctx, slug, percent)
}
