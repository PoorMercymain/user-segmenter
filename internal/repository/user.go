package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

type user struct {
	*postgres
}

func NewUser(pg *postgres) *user {
	return &user{pg}
}

func (r *user) UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error {
	log, err := logger.GetLogger()
	if err != nil {
		return err
	}

	conn, err := r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	slugs := append(slugsToAdd, slugsToDelete...)

	var str string

	for _, slug := range slugs {

		err = conn.QueryRow(ctx, "SELECT slug FROM slugs WHERE slug = $1", slug).Scan(&str)
		if err != nil {
			if err == pgx.ErrNoRows {
				return appErrors.ErrorNoRows
			}
			return err
		}
	}

	err = conn.QueryRow(ctx, "SELECT user_id FROM users WHERE user_id = $1", userID).Scan(&str)
	if err != nil {
		if err == pgx.ErrNoRows {
			conn.Exec(ctx, "INSERT INTO users VALUES ($1, $2)", userID, make([]string, 0))
		} else {
			return err
		}
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Infoln(err)
		}
	}()

	for _, slug := range slugsToAdd {
		updateResult, err := tx.Exec(ctx, "UPDATE users SET slugs = array_append(slugs, $1) WHERE user_id = $2 AND NOT $1 = ANY(slugs)", slug, userID)
		if err != nil {
			return err
		}

		if updateResult.RowsAffected() != 0 {
			_, err = tx.Exec(ctx, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", userID, slug, time.Now(), false)
			if err != nil {
				return err
			}
		}
	}

	for _, slug := range slugsToDelete {
		err = conn.QueryRow(ctx, "SELECT user_id FROM users WHERE user_id = $1 AND $2 = ANY(slugs)", userID, slug).Scan(&str)
		if err != nil {
			if err == pgx.ErrNoRows {
				return appErrors.ErrorNoRows
			}
			return err
		}
	}

	for _, slug := range slugsToDelete {
		updateResult, err := tx.Exec(ctx, "UPDATE users SET slugs = array_remove(slugs, $1) WHERE user_id = $2 AND $1 = ANY(slugs)", slug, userID)
		if err != nil {
			return err
		}
		if updateResult.RowsAffected() != 0 {
			_, err = tx.Exec(ctx, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", userID, slug, time.Now(), true)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *user) ReadUserSegments(ctx context.Context, userID string) ([]string, error) {
	conn, err := r.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var slugs []string

	err = conn.QueryRow(ctx, "SELECT slugs FROM users WHERE user_id = $1", userID).Scan(&slugs)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrorNoRows
		}
		return nil, err
	}

	return slugs, nil
}

func (r *user) CreateDeletionTime(ctx context.Context, userID string, slug string, deletionTime time.Time) error {
	log, err := logger.GetLogger()
	if err != nil {
		return err
	}

	conn, err := r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Infoln(err)
		}
	}()

	_, err = tx.Exec(ctx, "INSERT INTO deletion_times VALUES ($1, $2, $3) ON CONFLICT(user_id, slug) DO UPDATE SET deletion_timestamp = $3", userID, slug, deletionTime)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
