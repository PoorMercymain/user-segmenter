package domain

import "context"

type SegmentService interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error
	ReadUserSegments(ctx context.Context, userID string) ([]string, error)
}

//go:generate mockgen -destination=mocks/repo_mock.gen.go -package=mocks . SegmentRepository
type SegmentRepository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error
	ReadUserSegments(ctx context.Context, userID string) ([]string, error)
}
