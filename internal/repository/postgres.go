package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

type postgres struct {
	*pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *postgres {
	return &postgres{pool}
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
