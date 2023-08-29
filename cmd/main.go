package main

import (
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"

	"github.com/PoorMercymain/user-segmenter/internal/config"
	"github.com/PoorMercymain/user-segmenter/internal/handler"
	"github.com/PoorMercymain/user-segmenter/internal/middleware"
	"github.com/PoorMercymain/user-segmenter/internal/repository"
	"github.com/PoorMercymain/user-segmenter/internal/service"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

func init() {
	logger.InitLogger()
}

func router(pgPool *pgxpool.Pool) (*echo.Echo, error) {
	e := echo.New()

	segRep := repository.NewSegment(pgPool)
	segSrv := service.NewSegment(segRep)
	segHan := handler.NewSegment(segSrv)

	e.POST("/api/segment", segHan.CreateSegment, middleware.UseGzipReader())
	e.DELETE("/api/segment", segHan.DeleteSegment, middleware.UseGzipReader())
	e.POST("/api/user", segHan.UpdateUserSegments, middleware.UseGzipReader())
	e.GET("/api/user/:user", segHan.ReadUserSegments)
	e.GET("/api/user-history/:user", segHan.ReadUserSegmentsHistory)

	return e, nil
}

func main() {
	log, err := logger.GetLogger()
	if err != nil {
		return
	}
	log.Infoln("logger started")

	conf := config.GetServerConfig()

	pgPool, err := repository.ConnectToPostgres(conf.DatabaseURI)
	if err != nil {
		log.Infoln(err)
		return
	}

	defer pgPool.Close()

	r, err := router(pgPool)
	if err != nil {
		log.Infoln(err)
		return
	}

	if err = r.Start(strings.TrimPrefix(conf.ServerAddress, "http://")); err != nil {
		log.Infoln(err)
	}
}
