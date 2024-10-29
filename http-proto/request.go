package httpproto

import (
	"net/url"
)

type HttpRequest struct {
	RequestLine
	Headers Headers
	Query   url.Values
	Cookies CookieList
}

type RequestLine struct {
	Method  HttpMethod
	URI     string
	Version string
}
