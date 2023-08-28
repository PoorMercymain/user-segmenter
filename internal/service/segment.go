package service

import (
	"context"

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
