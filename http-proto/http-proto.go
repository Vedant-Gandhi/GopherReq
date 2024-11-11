package httpproto

import (
	"fmt"
	"http-v1_1/http-proto/common"
	"net"
	"os"
	"strings"
	"time"
)

const HEADER_LIMIT_BYTES = uint32(8192)

var supportedHttpMethods = []string{string(common.Get), string(common.Post), string(common.Put), string(common.Delete)}

type Config struct {
	Domain  string
	Timeout int
}

type HttpServer struct {
	listener net.Listener
	timeout  int
}

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
