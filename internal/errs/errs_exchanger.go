package errs

import "errors"

var (
	ErrNoRows               = errors.New("no exchange rates found")
	ErrUnsupportedInputCurr = errors.New("'from currency' not supported")
	ErrUnsupportedOutputCur = errors.New("'to currency' not supported")
)
