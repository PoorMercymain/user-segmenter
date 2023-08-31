package errors

import "errors"

var (
	ErrorInvalidRightLimit  = errors.New("invalid right limit, use more than 0")
	ErrorRightLimitIsTooLow = errors.New("right limit is too low, should be more or equal to amount")
)
