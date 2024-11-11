package httpproto

import (
	"errors"
	"http-v1_1/http-proto/common"
	"http-v1_1/http-proto/cookie"
	"io"
	"net/url"
	"slices"
	"strings"
	"unicode"
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

func parseRequestLine(rawData string) (reqLine RequestLine, err error) {
	// Format - Method SP Request-URI SP HTTP-Version
	// Ref - https://www.w3.org/Protocols/HTTP/1.1/draft-ietf-http-v11-spec-01#Request-Line

	rawData = strings.Trim(rawData, " ")

	indiviualData := strings.Split(rawData, " ")

	if len(indiviualData) != 3 {
		err = errors.New("invalid request line")
		return
	}

	rawMethod := indiviualData[0]
	rawMethod = strings.Trim(rawMethod, " ")

	if !slices.Contains(supportedHttpMethods, rawMethod) {
		err = errors.New("invalid http method")
		return
	}

	reqLine.Method = common.HttpMethod(rawMethod)

	// By default set the resource as self.
	selfResourceUrl, err := url.Parse("*")

	if err != nil {
		return reqLine, err
	}

	reqLine.URI = *selfResourceUrl

	if indiviualData[1] != "*" {
		uri, err := url.ParseRequestURI(indiviualData[1])

		if err != nil {
			return reqLine, err
		}

		reqLine.URI = *uri
	}

	reqLine.Version = indiviualData[2]

	return
}

func parseRequestHeaders(rawHeaders string) Headers {

	splitHeader := strings.Split(rawHeaders, "\r\n")

	headers := make(Headers)
	for _, line := range splitHeader {

		row := strings.SplitN(line, ":", 2)

		if len(row) != 2 {
			continue
		}

		// Trim spaces from key and value
		key := strings.Trim(row[0], " ")
		value := strings.Trim(row[1], " ")

		// Split key by "-" and capitalize each part
		parts := strings.Split(key, "-")
		for i, part := range parts {
			if len(part) > 0 {
				// Capitalize first letter and keep rest as is
				parts[i] = string(unicode.ToUpper(rune(part[0]))) + part[1:]
			}
		}

		// Join parts back with "-"
		finalKey := strings.Join(parts, "-")

		headers.Upsert(finalKey, HeaderValue(value))
	}
	return headers
}
