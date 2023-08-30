package main

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/PoorMercymain/user-segmenter/docs"
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

func router(pgPool *pgxpool.Pool, serverAddress string) (*echo.Echo, error) {
	e := echo.New()

	pg := repository.NewPostgres(pgPool)

	segRep := repository.NewSegment(pg)
	usrRep := repository.NewUser(pg)
	repRep := repository.NewReport(pg)

	segSrv := service.NewSegment(segRep)
	usrSrv := service.NewUser(usrRep)
	repSrv := service.NewReport(repRep)

	segHan := handler.NewSegment(segSrv)
	usrHan := handler.NewUser(usrSrv)
	repHan := handler.NewReport(repSrv)

	e.POST("/api/segment", segHan.CreateSegment, middleware.UseGzipReader())
	e.DELETE("/api/segment", segHan.DeleteSegment, middleware.UseGzipReader())
	e.POST("/api/user", usrHan.UpdateUserSegments, middleware.UseGzipReader())
	e.GET("/api/user/:user", usrHan.ReadUserSegments)
	e.GET("/api/user-history/:user", repHan.ReadUserSegmentsHistory, middleware.AddServerAddressToContext(serverAddress))
	e.GET("/api/reports/:report", repHan.ReadUserSegmentsHistoryReport)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log, err := logger.GetLogger()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			err := segRep.DeleteExpiredSegments(context.Background())
			if err != nil {
				log.Infoln(err)
			}
			time.Sleep(7 * time.Second)
		}
	}()

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

	r, err := router(pgPool, conf.ServerAddress)
	if err != nil {
		log.Infoln(err)
		return
	}

	if err = r.Start(strings.TrimPrefix(conf.ServerAddress, "http://")); err != nil {
		log.Infoln(err)
	}
}
