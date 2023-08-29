package domain

import (
	"context"
	"time"
)

type SegmentService interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error
	ReadUserSegments(ctx context.Context, userID string) ([]string, error)
	ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time) ([]HistoryElem, error)
	CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error
	AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error
}

//go:generate mockgen -destination=mocks/repo_mock.gen.go -package=mocks . SegmentRepository
type SegmentRepository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error
	ReadUserSegments(ctx context.Context, userID string) ([]string, error)
	ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time) ([]HistoryElem, error)
	CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error
	AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error
}
