package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"

	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

func TestUseGzipReader(t *testing.T) {
	logger.InitLogger()
	log, err := logger.GetLogger()
	require.NoError(t, err)

	e := echo.New()

	ts := httptest.NewServer(e)
	defer ts.Close()

	buf := bytes.NewBuffer([]byte(""))
	w := gzip.NewWriter(buf)
	w.Write([]byte("12345"))
	w.Close()
	r := bytes.NewBuffer(buf.Bytes())

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/test", r)
	require.NoError(t, err)

	req.Header.Set("Content-Encoding", "gzip")

	req2, err := http.NewRequest(http.MethodPost, ts.URL+"/test", r)
	require.NoError(t, err)

	e.POST("/test", func(ctx echo.Context) error {
		c := ctx.Request()
		if c.Header.Get("Content-Encoding") != "gzip" {
			log.Infoln("Content-Encoding is not gzip")
			ctx.Response().WriteHeader(http.StatusBadRequest)
			return nil
		} else {
			var b []byte
			b, err = io.ReadAll(ctx.Request().Body)
			ctx.Request().Body.Close()
			require.NoError(t, err)
			log.Infoln(string(b))
			ctx.Response().WriteHeader(http.StatusOK)
			return nil
		}
	}, UseGzipReader())

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	resp, err = ts.Client().Do(req2)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	req4, err := http.NewRequest(http.MethodPost, ts.URL+"/test", r)
	req4.Header.Set("Content-Encoding", "compress,deflate,br")
	require.NoError(t, err)

	resp, err = ts.Client().Do(req4)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	req5, err := http.NewRequest(http.MethodPost, ts.URL+"/test", nil)
	req5.Header.Set("Content-Encoding", "gzip")
	require.NoError(t, err)

	resp, err = ts.Client().Do(req5)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}
