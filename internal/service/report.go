package service

import (
	"context"
	"io"
	"time"

	"github.com/PoorMercymain/user-segmenter/internal/domain"
)

type report struct {
	repo domain.ReportRepository
}

func NewReport(repo domain.ReportRepository) *report {
	return &report{repo: repo}
}

func (s *report) ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time, limit, offset int) ([]domain.HistoryElem, error) {
	return s.repo.ReadUserSegmentsHistory(ctx, userID, startDate, endDate, limit, offset)
}

func (s *report) CreateCSV(ctx context.Context, userID string, startDate, endDate time.Time) (string, error) {
	return s.repo.CreateCSV(ctx, userID, startDate, endDate)
}

func (s *report) SendCSVReportFile(reportName string, writer io.Writer) error {
	return s.repo.SendCSVReportFile(reportName, writer)
}
