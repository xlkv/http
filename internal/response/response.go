package response

import (
	"fmt"
	"io"
	"strconv"

	"http.xlkv.io/internal/headers"
)

type StatusCode int

const (
	OK          StatusCode = 200
	BadRequest  StatusCode = 400
	ServerError StatusCode = 500
)

var StatusMessages = map[StatusCode]string{
	OK:          "HTTP/1.1 200 OK",
	BadRequest:  "HTTP/1.1 400 Bad Request",
	ServerError: "HTTP/1.1 500 Internal Server Error",
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		header := fmt.Sprintf("%v: %v\r\n", key, value)
		_, err := w.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return fmt.Errorf("headers is empty")
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusMessage, ok := StatusMessages[statusCode]
	if ok {
		statusLine := fmt.Sprintf("%v\r\n", statusMessage)
		w.Write([]byte(statusLine))
		return nil
	}
	return fmt.Errorf("Unknown status code")
}
