package httpproto

type HttpRequest struct {
	RequestLine
	Headers Headers
}

type RequestLine struct {
	Method  HttpMethod
	URI     string
	Version string
}
