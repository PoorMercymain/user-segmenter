package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

type report struct {
	*postgres
}

func NewReport(pg *postgres) *report {
	return &report{pg}
}

func (r *report) ReadUserSegmentsHistory(ctx context.Context, userID string, startDate, endDate time.Time, limit, offset int) ([]domain.HistoryElem, error) {
	conn, err := r.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var history []domain.HistoryElem

	var str string

	err = conn.QueryRow(ctx, "SELECT user_id FROM users WHERE user_id = $1", userID).Scan(&str)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrorNoRows
		}
		return nil, err
	}

	rows, err := conn.Query(ctx, "SELECT * FROM users_segment_history WHERE user_id = $1 AND modified_at <= $2 AND modified_at >= $3 ORDER BY modified_at DESC LIMIT $4 OFFSET $5", userID, endDate, startDate, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var historyElement domain.HistoryElem
		var isDeletion bool
		err = rows.Scan(&historyElement.UserID, &historyElement.Slug, &historyElement.DateTime, &isDeletion)
		if err != nil {
			return nil, err
		}

		historyElement.Operation = "addition"
		if isDeletion {
			historyElement.Operation = "deletion"
		}

		history = append(history, historyElement)
	}

	return history, nil
}

func (r *report) CreateCSV(ctx context.Context, userID string, startDate, endDate time.Time) (string, error) {
	filenamePattern := fmt.Sprintf("report*%d.csv", time.Now().UnixNano())

	log, _ := logger.GetLogger()

	f, err := os.CreateTemp("reports", filenamePattern)
	if err != nil {
		log.Infoln(err)
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'

	for i := 0; ; i++ {
		history, err := r.ReadUserSegmentsHistory(ctx, userID, startDate, endDate, 15, 15*i)
		if err != nil {
			return "", err
		}

		if len(history) == 0 {
			break
		}

		log.Infoln(len(history))

		for _, historyElement := range history {
			historyElementStrSlice := []string{historyElement.UserID, historyElement.Slug, historyElement.Operation, historyElement.DateTime.Format(time.RFC3339)}
			if err := w.Write(historyElementStrSlice); err != nil {
				return "", err
			}
		}
		w.Flush()
	}
	return f.Name(), nil
}

func (r *report) SendCSVReportFile(reportName string, writer io.Writer) error {
	log, _ := logger.GetLogger()
	_, err := os.Stat(".\\reports\\" + reportName)
	if err != nil {
		log.Infoln(err)
		if os.IsNotExist(err) {
			return appErrors.ErrorFileNotFound
		} else {
			return err
		}
	}

	f, err := os.Open(".\\reports\\" + reportName)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	written, err := io.CopyBuffer(writer, f, buffer)
	if err != nil {
		return err
	}

	if written == 0 {
		return appErrors.ErrorEmptyFile
	}

	return nil
}
