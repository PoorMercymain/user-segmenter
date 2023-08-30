package service

import (
	"context"

	"github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	slugvalidator "github.com/PoorMercymain/user-segmenter/pkg/slug-validator"
)

var (
	_ domain.SegmentService = (*segment)(nil)
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

func (s *segment) AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error {
	return s.repo.AddSegmentToPercentOfUsers(ctx, slug, percent)
}
