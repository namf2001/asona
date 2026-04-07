package pagination

const (
	Page    = 1
	Size    = 20
	MaxSize = 1000
)

// Input for pagination input
type Input struct {
	Page         int
	Size         int
	IncludeTotal bool
}

// Validate validates for pagination inputs
func (i Input) Validate() error {
	if i.Page <= 0 {
		return ErrPageMustBeGreaterThanZero
	}

	if i.Size <= 0 {
		return ErrSizeMustBeGreaterThanZero
	}

	if i.Size > MaxSize {
		return ErrExceededMaxPaginationSize
	}

	return nil
}
