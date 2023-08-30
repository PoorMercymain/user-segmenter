package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
	uniquenumbersgenerator "github.com/PoorMercymain/user-segmenter/pkg/unique-numbers-generator"
)

var (
	_ domain.SegmentRepository = (*segment)(nil)
)

type segment struct {
	*postgres
}

func NewSegment(pg *postgres) *segment {
	return &segment{pg}
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

	return tx.Commit(ctx)
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

	var sl string
	err = conn.QueryRow(ctx, "SELECT slug FROM slugs WHERE slug = $1", slug).Scan(&sl)
	if err == pgx.ErrNoRows {
		return appErrors.ErrorNoRows
	} else if err != nil {
		log.Infoln(err)
		return err
	}
	log.Infoln(sl)

	_, err = conn.Exec(ctx, "DELETE FROM slugs WHERE slug = $1", slug)
	if err != nil {
		return err
	}

	go func() {
		c := context.Background()

		conn, err := r.Acquire(c)
		if err != nil {
			log.Infoln(err)
			return
		}
		defer conn.Release()

		tx, err := conn.Begin(c)
		if err != nil {
			log.Infoln(err)
			return
		}
		defer func() {
			err = tx.Rollback(c)
			if err != nil {
				log.Infoln(err)
			}
		}()

		userIDs := make([]string, 0, 1)

		rows, err := tx.Query(c, "SELECT user_id FROM users WHERE $1 = ANY(slugs)", slug)
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
			_, err = tx.Exec(c, "UPDATE users SET slugs = array_remove(slugs, $1) WHERE user_id = $2", slug, id)
			if err != nil {
				log.Infoln(err)
			}
			_, err = tx.Exec(c, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", id, slug, time.Now(), true)
			if err != nil {
				log.Infoln(err)
			}
		}

		_, err = tx.Exec(c, "DELETE FROM deletion_times WHERE slug = $1", slug)
		if err != nil {
			log.Infoln(err)
		}
		err = tx.Commit(c)
		if err != nil {
			log.Infoln(err)
		}
	}()

	return nil
}

func (r *segment) AddSegmentToPercentOfUsers(ctx context.Context, slug string, percent int) error {
	log, err := logger.GetLogger()
	if err != nil {
		return err
	}

	conn, err := r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var usersAmount int
	err = conn.QueryRow(ctx, "SELECT COUNT(user_id) FROM users").Scan(&usersAmount)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Infoln(usersAmount)
			return appErrors.ErrorNoRows
		}
		return err
	}

	go func() {
		choosenAmount := int((float64(usersAmount) / 100) * float64(percent))
		log.Infoln(choosenAmount, usersAmount, percent)
		randomNumbersMap, err := uniquenumbersgenerator.GenerateUniqueNonNegativeNumbers(choosenAmount, usersAmount)
		if err != nil {
			log.Infoln(err)
			return
		}

		log.Infoln(randomNumbersMap)
		c := context.Background()

		conn, err := r.Acquire(c)
		if err != nil {
			log.Infoln(err)
			return
		}
		defer conn.Release()

		tx, err := conn.Begin(c)
		if err != nil {
			log.Infoln(err)
			return
		}
		defer func() {
			err = tx.Rollback(c)
			if err != nil {
				log.Infoln(err)
			}
		}()

		for randNum := range randomNumbersMap {
			var userID string
			err = tx.QueryRow(c, "SELECT user_id FROM users ORDER BY user_id DESC LIMIT 1 OFFSET $1 ", randNum).Scan(&userID)
			if err != nil {
				log.Infoln(err)
				return
			}
			_, err = tx.Exec(c, "UPDATE users SET slugs = array_append(slugs, $1) WHERE user_id = $2", slug, userID)
			if err != nil {
				log.Infoln(err)
				return
			}
			_, err = tx.Exec(c, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", userID, slug, time.Now(), false)
			if err != nil {
				log.Infoln(err)
				return
			}
		}
		err = tx.Commit(c)
		if err != nil {
			log.Infoln(err)
		}
	}()

	return nil
}

func (r *segment) DeleteExpiredSegments(ctx context.Context) error {
	log, err := logger.GetLogger()
	if err != nil {
		return err
	}

	conn, err := r.Acquire(ctx)
	if err != nil {
		log.Infoln(err)
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT user_id, slug FROM deletion_times WHERE deletion_timestamp <= $1", time.Now())
	if err != nil {
		return err
	}

	for rows.Next() {
		var userID, slug string
		err = rows.Scan(&userID, &slug)
		if err != nil {
			return err
		}

		conn, err := r.Acquire(ctx)
		if err != nil {
			log.Infoln(err)
			return err
		}
		defer conn.Release()

		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Infoln(err)
			return err
		}
		defer func() {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Infoln(err)
			}
		}()

		execResult, err := tx.Exec(ctx, "UPDATE users SET slugs = array_remove(slugs, $1) WHERE user_id = $2 AND $1 = ANY(slugs)", slug, userID)
		if err != nil {
			return err
		}
		t := time.Now()
		_, err = tx.Exec(ctx, "DELETE FROM deletion_times WHERE user_id = $1 AND slug = $2", userID, slug)
		if err != nil {
			return err
		}
		if execResult.RowsAffected() != 0 {
			_, err = tx.Exec(ctx, "INSERT INTO users_segment_history VALUES ($1, $2, $3, $4)", userID, slug, t, true)
			if err != nil {
				log.Infoln(err)
				return err
			}
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Infoln(err)
			return err
		}
	}
	return nil
}
