package pagination

import (
	"errors"
)

var (
	ErrPageMustBeGreaterThanZero = errors.New("page must be greater than 0")
	ErrSizeMustBeGreaterThanZero = errors.New("size must be greater than 0")
	ErrExceededMaxPaginationSize = errors.New("maximum size for each page has been exceeded")
)
