package middleware

import (
	"compress/gzip"
	"net/http"

	"github.com/labstack/echo"

	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

func UseGzipReader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log, err := logger.GetLogger()
			if err != nil {
				c.Response().WriteHeader(http.StatusInternalServerError)
				return err
			}
			log.Infoln("in gzip")
			if len(c.Request().Header.Values("Content-Encoding")) == 0 {
				log.Infoln("no gzip")
				return next(c)
			}
			for i, headerValue := range c.Request().Header.Values("Content-Encoding") {
				if headerValue == "gzip" {
					break
				}
				log.Infoln(i, (len(c.Request().Header.Values("Content-Encoding")) - 1))
				if i == (len(c.Request().Header.Values("Content-Encoding")) - 1) {
					log.Infoln("no gzip")
					return next(c)
				}
			}

			gzipReader, err := gzip.NewReader(c.Request().Body)
			if err != nil {
				log.Infoln("no")
				c.Response().WriteHeader(http.StatusBadRequest)
				return err
			}
			c.Request().Body.Close()

			c.Request().Body = gzipReader

			return next(c)
		}
	}
}
