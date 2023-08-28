package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

type segment struct {
	*pgxpool.Pool
}

func NewSegment(pool *pgxpool.Pool) *segment {
	return &segment{pool}
}

func ConnectToPostgres(DSN string) (*pgxpool.Pool, error) {
	log, err := logger.GetLogger()
	if err != nil {
		return nil, err
	}

	config, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		log.Infoln(err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Infoln(err)
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *segment) CreateSegment(ctx context.Context, slug string) error {
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

	var pgErr *pgconn.PgError
	_, err = tx.Exec(ctx, "INSERT INTO slugs VALUES ($1)", slug)
	errors.As(err, &pgErr)
	if err != nil && pgErr.Code == pgerrcode.UniqueViolation {
		return appErrors.ErrorUniqueViolation
	} else if err != nil {
		return err
	}

	return nil
}

func (r *segment) DeleteSegment(ctx context.Context, slug string) error {
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

	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM slugs WHERE slug = $1", slug).Scan()
	if err == pgx.ErrNoRows {
		return appErrors.ErrorNoRows
	} else if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM slugs WHERE slug = $1", slug)
	if err != nil {
		return err
	}

	go func() {
		userIDs := make([]string, 0, 1)

		rows, err := tx.Query(ctx, "SELECT user_id FROM users WHERE $1 = ANY(slugs)", slug)
		if err != nil {
			log.Infoln(err)
		}

		for rows.Next() {
			var id string
			err = rows.Scan(&id)
			if err != nil {
				log.Infoln(err)
				break
			}
			userIDs = append(userIDs, id)
		}

		for _, id := range userIDs {
			_, err = tx.Exec(ctx, "UPDATE users SET slugs = array_remove(slugs, $1) WHERE user_id = $2", slug, id)
			if err != nil {
				log.Infoln(err)
			}
			_, err = tx.Exec(ctx, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", id, slug, time.Now(), true)
			if err != nil {
				log.Infoln(err)
			}
		}

		_, err = tx.Exec(ctx, "DELETE FROM deletion_times WHERE slug = $1", slug)
		if err != nil {
			log.Infoln(err)
		}

	}()

	return nil
}

func (r *segment) UpdateUserSegments(ctx context.Context, userID string, slugsToAdd []string, slugsToDelete []string) error {
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

	for _, slug := range slugs {
		err = conn.QueryRow(ctx, "SELECT slug FROM slugs WHERE slug = $1", slug).Scan()
		if err != nil {
			if err == pgx.ErrNoRows {
				return appErrors.ErrorNoRows
			}
			return err
		}
	}

	err = conn.QueryRow(ctx, "SELECT user_id FROM users WHERE user_id = $1", userID).Scan()
	if err != nil {
		if err == pgx.ErrNoRows {
			conn.Exec(ctx, "INSERT INTO users VALUES ($1, $2)", userID, make([]string, 0))
		} else {
			return err
		}
	}

	for _, slug := range slugsToDelete {
		err = conn.QueryRow(ctx, "SELECT user_id FROM users WHERE user_id = $1 AND $2 = ANY(slugs)", userID, slug).Scan()
		if err != nil {
			if err == pgx.ErrNoRows {
				return appErrors.ErrorNoRows
			}
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

	for _, slug := range slugsToDelete {
		_, err = tx.Exec(ctx, "UPDATE users SET slugs = array_remove(slugs, $1) WHERE user_id = $2", slug, userID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", userID, slug, time.Now(), true)
		if err != nil {
			return err
		}
	}

	for _, slug := range slugsToAdd {
		_, err = tx.Exec(ctx, "UPDATE users SET slugs = array_append(slugs, $1) WHERE user_id = $2", slug, userID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", userID, slug, time.Now(), false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *segment) ReadUserSegments(ctx context.Context, userID string) ([]string, error) {
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
