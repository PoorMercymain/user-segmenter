package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/PoorMercymain/user-segmenter/internal/domain"
)

func TestAddServerAddressToContext(t *testing.T) {
	e := echo.New()

	ts := httptest.NewServer(e)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/test", bytes.NewReader([]byte("")))
	require.NoError(t, err)

	req1, err := http.NewRequest(http.MethodGet, ts.URL+"/test1", bytes.NewReader([]byte("")))
	require.NoError(t, err)

	e.GET("/test", func(ctx echo.Context) error {
		addr := ctx.Request().Context().Value(domain.Key("server"))
		require.Equal(t, "test", addr)
		return nil
	}, AddServerAddressToContext("test"))
	e.GET("/test1", func(ctx echo.Context) error {
		addr := ctx.Request().Context().Value(domain.Key("server"))
		require.Nil(t, addr)
		return nil
	})

	resp, err := ts.Client().Do(req)
	resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = ts.Client().Do(req1)
	resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
