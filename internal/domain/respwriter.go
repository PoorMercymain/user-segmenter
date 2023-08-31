package domain

import (
	"io"
	"net/http"
)

type RespWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w RespWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
