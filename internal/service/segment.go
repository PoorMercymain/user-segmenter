package service

import (
	"context"
	"strings"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	slugValidator "github.com/PoorMercymain/user-segmenter/pkg/slug-validator"
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
	if slugValidator.IsSlug(strings.ToLower(slug)) {
		return s.repo.CreateSegment(ctx, slug)
	}

	return appErrors.ErrorNotASlug
}

func (s *segment) DeleteSegment(ctx context.Context, slug string) error {
	return s.repo.DeleteSegment(ctx, slug)
}

func (s *segment) AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error {
	return s.repo.AddSegmentToPercentOfUsers(ctx, slug, percent)
}
