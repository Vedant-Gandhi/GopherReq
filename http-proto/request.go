package httpproto

import (
	"http-v1_1/http-proto/cookie"
	"io"
	"net/url"
)

type HttpRequest struct {
	RequestLine
	Headers Headers
	Query   url.Values
	Cookies cookie.CookieList
	Body    RequestBody
}

type RequestLine struct {
	Method  HttpMethod
	URI     string
	Version string
}

type RequestBody io.Reader
