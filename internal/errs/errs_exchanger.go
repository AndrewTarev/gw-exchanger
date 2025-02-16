package errs

import "errors"

var (
	ErrNoRows = errors.New("no exchange rates found")
)
