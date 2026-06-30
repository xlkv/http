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

type writeState int

const (
	stateStatusLine writeState = iota
	stateHeaders
	stateBody
)

type Writer struct {
	writer io.Writer
	state  writeState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  stateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateStatusLine {
		return fmt.Errorf("state not init.")
	}
	statusMessage, ok := StatusMessages[statusCode]
	if ok {
		statusLine := fmt.Sprintf("%v\r\n", statusMessage)
		w.writer.Write([]byte(statusLine))
		w.state = stateHeaders
		return nil
	}
	return fmt.Errorf("Unknown status code")
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateHeaders {
		return fmt.Errorf("state error")
	}
	for key, value := range headers {
		header := fmt.Sprintf("%v: %v\r\n", key, value)
		_, err := w.writer.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	w.state = stateBody
	if err != nil {
		return err
	}
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("state error")
	}
	n, err := w.writer.Write(p)

	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("state error")
	}
	chunkSize := fmt.Sprintf("%x\r\n", len(p))
	w.WriteBody([]byte(chunkSize))
	n, err := w.WriteBody([]byte(p))
	w.WriteBody([]byte("\r\n"))
	return n, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("state error")
	}
	return w.WriteBody([]byte("0\r\n\r\n"))
}

var StatusMessages = map[StatusCode]string{
	OK:          "HTTP/1.1 200 OK",
	BadRequest:  "HTTP/1.1 400 Bad Request",
	ServerError: "HTTP/1.1 500 Internal Server Error",
}

func (w *Writer) GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"content-length": strconv.Itoa(contentLen),
		"connection":     "close",
		"content-type":   "text/plain",
	}
}
