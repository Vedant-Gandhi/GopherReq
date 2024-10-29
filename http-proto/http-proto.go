package httpproto

import (
	"bytes"
	"errors"
	"fmt"
	"http-v1_1/http-proto/common"
	"io"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"
	"unicode"
)

type HttpMethod string

const (
	Get    HttpMethod = "GET"
	Post   HttpMethod = "POST"
	Put    HttpMethod = "PUT"
	Delete HttpMethod = "DELETE"
)

const HEADER_LIMIT_BYTES = 8192

var supportedHttpMethods = []string{string(Get), string(Post), string(Put), string(Delete)}

type Headers map[string]string

type Config struct {
	Domain  string
	Timeout int
}

type HttpServer struct {
	listener net.Listener
}

func NewServer(cfg Config) (server HttpServer, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:8811")
	if err != nil {
		fmt.Printf("Error while listening to the socket: %v\n", err)
		return
	}

	server.listener = listener
	return
}

func (s *HttpServer) Listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error while accepting the connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func (s *HttpServer) ShutDown() {
	s.listener.Close()
}

func handleConnection(conn net.Conn) {

	request, _ := readHeader(conn)
	request, _ = parseQuery(request)

	response := generateHttpResponse(request)

	encodedResponse := encodeResponseforTransfer(response)

	conn.Write(encodedResponse)

	defer conn.Close()

}

func readHeader(conn net.Conn) (request HttpRequest, err error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	data := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	// This loop keeps on reading headers if it does not fit in one buffer.
	for {
		bytesReadCount, err := conn.Read(readBuffer)
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

		if data.Len() > HEADER_LIMIT_BYTES {
			break
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

			parsedHeaders := parseHeadertoMap(headers[reqLineIdx+2:]) // Added +2 to skip \r\n

			request.Headers = parsedHeaders
			request.RequestLine = parsedReqLine

			if len(request.Headers["Cookie"]) != 0 {
				request.Cookies = NewCookieList()
				splitCookies := strings.Split(request.Headers["Cookie"], ";")
				for _, cookie := range splitCookies {
					c, err := parseRequestCookie(cookie)
					if err != nil {
						fmt.Printf("Error cookie is not valid - %v", err)
						continue
					}

					request.Cookies.Add(c)

				}
			}

			return request, nil // Return immediately after parsing headers
		}
	}

	return request, errors.New("incomplete headers")
}

func parseRequestLine(rawData string) (reqLine RequestLine, err error) {
	// Format - Method SP Request-URI SP HTTP-Version CRLF
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

	reqLine.Method = HttpMethod(rawMethod)

	// By default set the resource as self.
	reqLine.URI = "*"

	if indiviualData[1] != "*" {
		uri, err := url.ParseRequestURI(indiviualData[1])

		if err != nil {
			return reqLine, err
		}

		reqLine.URI = uri.String()
	}

	reqLine.Version = indiviualData[2]

	return

}

func parseHeadertoMap(rawHeaders string) Headers {

	splitHeader := strings.Split(rawHeaders, "\r\n")

	headerMap := make(map[string]string)
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

		headerMap[finalKey] = value
	}
	return headerMap
}

func generateHttpResponse(request HttpRequest) (response HttpResponse) {

	respCode := common.StatusCode(OK)

	if request.Method != Get {
		respCode = common.StatusCode(NOT_IMPLEMENTED)
	}

	respLine := ResponseLine{
		Code:    respCode,
		Reason:  httpStatusPhraseReasons[respCode],
		Version: request.Version,
	}

	headers := Headers{
		"Date":           time.Now().UTC().Format(time.RFC1123),
		"Content-Length": "0",
	}

	response = HttpResponse{
		ResponseLine: respLine,
		Headers:      headers,
	}

	return

}

func encodeResponseforTransfer(response HttpResponse) (data []byte) {

	wireResponse := strings.Builder{}

	wireResponse.WriteString(fmt.Sprintf("%s %d %s%s", response.ResponseLine.Version, response.ResponseLine.Code, response.ResponseLine.Reason, common.CRLF)) // The Response Line.

	for key, value := range response.Headers {
		wireResponse.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	wireResponse.WriteString(common.CRLF)

	data = []byte(wireResponse.String())
	return
}

func parseQuery(request HttpRequest) (HttpRequest, error) {

	if len(request.URI) == 0 {
		return request, nil
	}

	query, err := url.ParseQuery(request.URI)

	if err != nil {
		fmt.Printf("Error while parsing the query variables - %v", err)
		return request, err
	}

	request.Query = query

	return request, nil
}

func parseCookies(request HttpRequest) (HttpRequest, error) {

	if len(request.URI) == 0 {
		return request, nil
	}

	query, err := url.ParseQuery(request.URI)

	if err != nil {
		fmt.Printf("Error while parsing the query variables - %v", err)
		return request, err
	}

	request.Query = query

	return request, nil
}

/** Can be used for future Post requests.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New client connected!")

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	data := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)
	headerComplete := false
	var contentLength int
	var bodyStart int

	for {
		readBytes, err := conn.Read(readBuffer)
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
			return
		}

		if readBytes == 0 {
			break
		}

		data.Write(readBuffer[:readBytes])

		// If we haven't found the end of headers yet, look for it
		if !headerComplete {
			if headerEnd := bytes.Index(data.Bytes(), []byte("\r\n\r\n")); headerEnd != -1 {
				headerComplete = true
				headers := string(data.Bytes()[:headerEnd])

				// Look for Content-Length header
				for _, line := range strings.Split(headers, "\r\n") {
					if strings.HasPrefix(strings.ToLower(line), "content-length:") {
						fmt.Sscanf(line, "Content-Length: %d", &contentLength)
						bodyStart = headerEnd + 4 // Skip the \r\n\r\n
						break
					}
				}

				// If no content length or content length is 0, we're done
				if contentLength == 0 {
					fmt.Println("Request complete - no body expected")
					break
				}
			}
		}

		// If we have found headers and have a content length, check if we've received the full body
		if headerComplete && contentLength > 0 {
			if data.Len()-bodyStart >= contentLength {
				fmt.Println("Request complete - received full body")
				break
			}
		}
	}

	fmt.Printf("Final buffer contents (%d bytes):\n%s\n", data.Len(), data.String())
}

*/
