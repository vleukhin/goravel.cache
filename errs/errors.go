package errs

import "errors"

var (
	ErrCacheMiss = errors.New(("cache: missed"))
)
