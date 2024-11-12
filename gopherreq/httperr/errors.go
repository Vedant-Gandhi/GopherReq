package httperr

import "errors"

// Shared Errors
var (
	ErrIncompleteHeader    = errors.New("incomplete headers")
	ErrHeaderLimitExceeded = errors.New("size of headers exceeds the limit")
	ErrInvalidHttpMethod   = errors.New("invalid http method")
)

// Http Request Errors
var (
	ErrInvalidContentLength = errors.New("content length is invalid")
	ErrInvalidRequestLine   = errors.New("invalid request line")
)
