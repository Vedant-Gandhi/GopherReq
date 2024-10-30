package httpproto

import (
	"http-v1_1/http-proto/cookie"
	"net/url"
)

type HttpRequest struct {
	RequestLine
	Headers Headers
	Query   url.Values
	Cookies cookie.CookieList
}

type RequestLine struct {
	Method  HttpMethod
	URI     string
	Version string
}
