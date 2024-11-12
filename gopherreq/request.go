package gopherreq

import (
	"bytes"
	"fmt"
	"gopherreq/gopherreq/common"
	"gopherreq/gopherreq/cookie"
	"gopherreq/gopherreq/httperr"

	"io"
	"net"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
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
		err = httperr.ErrInvalidRequestLine
		return
	}

	rawMethod := indiviualData[0]
	rawMethod = strings.Trim(rawMethod, " ")

	if !slices.Contains(supportedHttpMethods, common.HttpMethod(rawMethod)) {
		err = httperr.ErrInvalidHttpMethod
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

		headers.Apsert(finalKey, HeaderValue(value))
	}
	return headers
}

// Reads the header from the connection.
func (h HttpServer) readHeader(conn net.Conn) (request HttpRequest, err error) {

	// Adjust the read deadline.
	conn.SetReadDeadline(time.Now().Add(time.Duration(h.timeout) * time.Millisecond))

	data := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	// This loop keeps on reading headers if it does not fit in one buffer.
	for {
		bytesReadCount, err := conn.Read(readBuffer)

		// Update the read deadline.
		conn.SetReadDeadline(time.Now().Add(time.Duration(h.timeout * int(time.Millisecond))))

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection")
				break
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Read timeout - no more data expected")
				break
			}
			fmt.Printf("Error while reading from the connection: %v\n", err)
			return request, err
		}

		data.Write(readBuffer[:bytesReadCount])

		if uint32(data.Len()) > HEADER_LIMIT_BYTES {
			fmt.Printf("Header len limit: %v", data.Len())
			err = httperr.ErrHeaderLimitExceeded
			return request, err
		}

		headerEnd := bytes.Index(data.Bytes(), []byte("\r\n\r\n"))
		if headerEnd != -1 { // If we have found header end
			headers := string(data.Bytes()[:headerEnd])
			reqLineIdx := strings.Index(headers, "\r\n")

			reqLine := headers[:reqLineIdx]
			parsedReqLine, err := parseRequestLine(reqLine)
			if err != nil {
				return request, err
			}

			parsedHeaders := parseRequestHeaders(headers[reqLineIdx+2:]) // Added +2 to skip \r\n

			request.Headers = parsedHeaders

			host := request.Headers.Get("host")

			if host != "" {
				parsedReqLine.URI.Host = host.String()
				request.URI = parsedReqLine.URI
				request.Method = parsedReqLine.Method
				request.Version = parsedReqLine.Version
				request.RawURI = parsedReqLine.URI.String()
			}
			return request, nil // Return immediately after parsing headers
		}
	}

	return request, httperr.ErrIncompleteHeader
}

func parseRequestCookie(request *HttpRequest) error {

	if len(request.Headers["Cookie"]) != 0 {

		request.Cookies = cookie.NewCookieList()

		cookieValues := request.Headers.GetAllValues("Cookie")

		for _, rawCookie := range cookieValues {

			splitCookies := strings.Split(rawCookie.String(), ";")

			for _, splitCookie := range splitCookies {

				c, err := cookie.ParseRequestCookie(splitCookie)
				if err != nil {
					fmt.Printf("Error cookie is not valid - %v", err)
					continue
				}
				request.Cookies.Add(c)

			}
		}

	}

	return nil
}

/**
 * This function reads the body from the request and stores in binary form.
 */
func (req *HttpRequest) readBody(conn net.Conn) (err error) {

	rawLen := "0"

	contentLength := req.Headers.Get("Content-Length")

	if contentLength != "" {
		rawLen = contentLength.String()
	}
	bodyLen, err := strconv.Atoi(rawLen)

	if err != nil {
		err = httperr.ErrInvalidContentLength
		return
	}

	if bodyLen == 0 {
		return
	}

	// Read the whole body at once in the buffer.
	buffer := make([]byte, bodyLen)
	_, err = conn.Read(buffer)
	req.Body = bytes.NewReader(buffer)

	return err
}
