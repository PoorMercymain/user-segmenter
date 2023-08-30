package errors

import "errors"

var (
	ErrorFileNotFound = errors.New("file not found")
	ErrorEmptyFile    = errors.New("the file is empty")
)
