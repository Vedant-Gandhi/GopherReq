package common

type StatusCode int

const CRLF = "\r\n"

type HttpMethod string

const (
	Get    HttpMethod = "GET"
	Post   HttpMethod = "POST"
	Put    HttpMethod = "PUT"
	Delete HttpMethod = "DELETE"
)
