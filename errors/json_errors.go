package errors

import "errors"

var (
	ErrorDuplicateInJSON = errors.New("duplicate key found in JSON")
)
