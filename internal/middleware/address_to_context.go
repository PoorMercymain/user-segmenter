package middleware

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/PoorMercymain/user-segmenter/internal/domain"
)

func AddServerAddressToContext(serverAddress string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), domain.Key("server"), serverAddress)))

			return next(c)
		}
	}
}
