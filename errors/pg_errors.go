package errors

import "errors"

var (
	ErrorUniqueViolation = errors.New("unique violation error")
	ErrorNoRows          = errors.New("no rows found for the provided slug")
)
