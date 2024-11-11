package httperr

import "errors"

// Shared Errors
var (
	ErrIncompleteHeader  = errors.New("incomplete headers")
	ErrInvalidHttpMethod = errors.New("invalid http method")
)

// Http Request Errors
var (
	ErrInvalidContentLength = errors.New("content length is invalid")
	ErrInvalidRequestLine   = errors.New("invalid request line")
)
