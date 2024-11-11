package httperr

import "errors"

var (
	ErrIncompleteHeader = errors.New("incomplete headers")
)
