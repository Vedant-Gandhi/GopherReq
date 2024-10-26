package httpproto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"
)

// HTTP method constants
type HTTPMethod string

const (
	MethodGet    HTTPMethod = "GET"
	MethodPost   HTTPMethod = "POST"
	MethodPut    HTTPMethod = "PUT"
	MethodDelete HTTPMethod = "DELETE"
)

var (
	supportedHTTPMethods = []string{
		string(MethodGet),
		string(MethodPost),
		string(MethodPut),
		string(MethodDelete),
	}

	ErrInvalidRequestLine = errors.New("invalid request line")
	ErrInvalidHTTPMethod  = errors.New("invalid http method")
	ErrIncompleteHeaders  = errors.New("incomplete headers")
)

// Headers represents HTTP headers as a map
type Headers map[string]string

// RequestLine represents the first line of an HTTP request
type RequestLine struct {
	Method  HTTPMethod
	URI     url.URL
	Version string
}

// HTTPRequest represents a complete HTTP request
type HTTPRequest struct {
	RequestLine
	Headers Headers // Make it public if needed by other packages
}

// ServerConfig holds configuration for the HTTP server
type ServerConfig struct {
	Address string        // Address to listen on
	Timeout time.Duration // Read timeout duration
}

// DefaultConfig returns default server configuration
func DefaultConfig() ServerConfig {
	return ServerConfig{
		Address: "127.0.0.1:8811",
		Timeout: 10 * time.Second,
	}
}

// Server represents an HTTP server
type Server struct {
	config   ServerConfig
	listener net.Listener
}

// NewServer creates a new HTTP server with the given configuration
func NewServer(cfg ServerConfig) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &Server{
		config:   cfg,
		listener: listener,
	}, nil
}

// Listen starts accepting connections
func (s *Server) Listen() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	return s.listener.Close()
}

// handleConnection processes a single connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	request, err := s.readRequest(conn)
	if err != nil {
		// Log error and possibly send error response
		fmt.Printf("Error reading request: %v\n", err)
		return
	}

	// Handle the request (to be implemented)
	fmt.Printf("Received request: %+v\n", request)
}

// readRequest reads and parses an HTTP request from a connection
func (s *Server) readRequest(conn net.Conn) (*HTTPRequest, error) {
	conn.SetReadDeadline(time.Now().Add(s.config.Timeout))

	data := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	for {
		n, err := conn.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			return nil, fmt.Errorf("failed to read from connection: %w", err)
		}

		data.Write(readBuffer[:n])

		if request, complete := s.parseRequest(data.Bytes()); complete {
			return request, nil
		}
	}

	return nil, ErrIncompleteHeaders
}

// parseRequest attempts to parse an HTTP request from raw bytes
func (s *Server) parseRequest(data []byte) (*HTTPRequest, bool) {
	headerEnd := bytes.Index(data, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return nil, false
	}

	headers := string(data[:headerEnd])
	reqLineEnd := strings.Index(headers, "\r\n")
	if reqLineEnd == -1 {
		return nil, false
	}

	reqLine, err := parseRequestLine(headers[:reqLineEnd])
	if err != nil {
		return nil, false
	}

	parsedHeaders := parseHeaders(headers[reqLineEnd+2:])

	return &HTTPRequest{
		RequestLine: reqLine,
		Headers:     parsedHeaders,
	}, true
}

// parseRequestLine parses the request line of an HTTP request
func parseRequestLine(raw string) (RequestLine, error) {
	parts := strings.Fields(raw)
	if len(parts) != 3 {
		return RequestLine{}, ErrInvalidRequestLine
	}

	if !slices.Contains(supportedHTTPMethods, parts[0]) {
		return RequestLine{}, ErrInvalidHTTPMethod
	}

	uri, err := url.ParseRequestURI(parts[1])
	if err != nil {
		return RequestLine{}, fmt.Errorf("invalid URI: %w", err)
	}

	return RequestLine{
		Method:  HTTPMethod(parts[0]),
		URI:     *uri,
		Version: parts[2],
	}, nil
}

// parseHeaders parses HTTP headers from a string
func parseHeaders(raw string) Headers {
	headers := make(Headers)
	lines := strings.Split(raw, "\r\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := normalizeHeaderKey(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	return headers
}

// normalizeHeaderKey normalizes HTTP header keys
func normalizeHeaderKey(key string) string {
	parts := strings.Split(key, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "-")
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
