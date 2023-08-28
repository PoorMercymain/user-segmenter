package errors

import "errors"

var (
	ErrorLoggerNotInitialized = errors.New("logger should be initialized, but it is nil")
)
