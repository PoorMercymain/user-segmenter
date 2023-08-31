package domain

import (
	"context"
	"io"
	"time"
)

type SegmentService interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error
}

type UserService interface {
	UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error
	ReadUserSegments(ctx context.Context, userID string) ([]string, error)
	CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error
}

type ReportService interface {
	ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time, limit, offset int) ([]HistoryElem, error)
	CreateCSV(ctx context.Context, userID string, startDate, endDate time.Time) (string, error)
	SendCSVReportFile(reportName string, writer io.Writer) error
}

//go:generate mockgen -destination=mocks/segment_repo_mock.gen.go -package=mocks . SegmentRepository
type SegmentRepository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error
	DeleteExpiredSegments(ctx context.Context) error
}

//go:generate mockgen -destination=mocks/user_repo_mock.gen.go -package=mocks . UserRepository
type UserRepository interface {
	UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error
	ReadUserSegments(ctx context.Context, userID string) ([]string, error)
	CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error
}

//go:generate mockgen -destination=mocks/report_repo_mock.gen.go -package=mocks . ReportRepository
type ReportRepository interface {
	ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time, limit, offset int) ([]HistoryElem, error)
	CreateCSV(ctx context.Context, userID string, startDate, endDate time.Time) (string, error)
	SendCSVReportFile(reportName string, writer io.Writer) error
}
