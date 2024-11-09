package httpproto

import (
	"http-v1_1/http-proto/common"
	"http-v1_1/http-proto/cookie"
	"io"
	"net/url"
)

type HttpRequest struct {
	Headers Headers           // The headers received from the client.
	Cookies cookie.CookieList // Stores the cookies received by the client in parsed format. These are cleaned and stored.
	Body    RequestBody       // The request body received from the client.
	Method  common.HttpMethod // The HTTP method for the request.
	URI     url.URL           // The URI for the request. It is parsed and clean version. You can read the query variables from here.
	Version string            // The HTTP Version for the request.
	RawURI  string            // The raw unformatted version of the uri as received from the client. Always use URI wherever possible instead of this.It is not sanitized and may lead to attacks.
}

type RequestBody io.Reader
