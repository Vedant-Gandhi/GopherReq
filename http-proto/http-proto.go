package httpproto

import (
	"bytes"
	"errors"
	"fmt"
	"http-v1_1/http-proto/common"
	"http-v1_1/http-proto/cookie"
	"io"
	"net"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const HEADER_LIMIT_BYTES = 8192

var supportedHttpMethods = []string{string(common.Get), string(common.Post), string(common.Put), string(common.Delete)}

type Config struct {
	Domain  string
	Timeout int
}

type HttpServer struct {
	listener net.Listener
	timeout  int
}

var (
	ErrInvalidContentLength = errors.New("content length is invalid")
)

func NewServer(cfg Config) (server HttpServer, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:8811")
	if err != nil {
		fmt.Printf("Error while listening to the socket: %v\n", err)
		return
	}

	server.listener = listener
	server.timeout = cfg.Timeout

	return
}

func (s *HttpServer) Listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error while accepting the connection: %v\n", err)
			continue
		}

		conn.SetDeadline(time.Now().Add(time.Duration(s.timeout) * time.Millisecond))

		go s.handleConnection(conn)
	}
}

func (s *HttpServer) ShutDown() {
	s.listener.Close()
}

func (s *HttpServer) handleConnection(conn net.Conn) {

	request, err := s.readHeader(conn)
	if err != nil {
		fmt.Printf("error while reading the header %v:", err)
		os.Exit(1)
	}

	err = parseRequestCookie(&request)
	if err != nil {
		fmt.Printf("error while reading the cookies %v:", err)
		os.Exit(1)
	}

	err = request.readBody(conn)
	if err != nil {
		fmt.Printf("error while reading the body %v:", err)
		os.Exit(1)
	}

	response := generateHttpWireResponse(request)

	encodedResponse := encodeHttpWireResponseToBinary(response)

	conn.Write(encodedResponse)

	defer conn.Close()

}

func (h HttpServer) readHeader(conn net.Conn) (request HttpRequest, err error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	data := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	// This loop keeps on reading headers if it does not fit in one buffer.
	for {
		bytesReadCount, err := conn.Read(readBuffer)

		// Update the read deadline.
		if h.timeout != 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(h.timeout * int(time.Millisecond))))
		}
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
			fmt.Printf("Header len limit: %v", data.Len())
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

			parsedHeaders := parseRequestHeaders(headers[reqLineIdx+2:]) // Added +2 to skip \r\n

			request.Headers = parsedHeaders

			host, hostExist := request.Headers.Get("host")

			if hostExist {
				parsedReqLine.URI.Host = host.String()
				request.URI = parsedReqLine.URI
				request.Method = parsedReqLine.Method
				request.Version = parsedReqLine.Version
				request.RawURI = parsedReqLine.URI.String()

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

func generateHttpWireResponse(request HttpRequest) (response HttpWireResponse) {

	respCode := common.StatusCode(OK)

	if request.Method != common.Get {
		respCode = common.StatusCode(NOT_IMPLEMENTED)
	}

	respLine := ResponseLine{
		Code:    respCode,
		Reason:  httpStatusPhraseReasons[respCode],
		Version: request.Version,
	}

	headers := make(Headers)

	headers.Set("Date", HeaderValue(time.Now().UTC().Format(time.RFC1123)))
	headers.Set("Content-Length", "0")

	response = HttpWireResponse{
		ResponseLine: respLine,
		Headers:      headers,
	}

	return

}

func encodeHttpWireResponseToBinary(response HttpWireResponse) (data []byte) {

	binaryResponse := strings.Builder{}

	binaryResponse.WriteString(fmt.Sprintf("%s %d %s%s", response.ResponseLine.Version, response.ResponseLine.Code, response.ResponseLine.Reason, common.CRLF)) // The Response Line.

	for key, value := range response.Headers {
		binaryResponse.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	binaryResponse.WriteString(common.CRLF)

	data = []byte(binaryResponse.String())
	return
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

	contentLength, exist := req.Headers.Get("Content-Length")

	if exist {
		rawLen = contentLength.String()
	}
	bodyLen, err := strconv.Atoi(rawLen)

	if err != nil {
		err = ErrInvalidContentLength
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
