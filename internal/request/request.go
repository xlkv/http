package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type parserState int

const (
	requestStateInitialized parserState = iota
	requestStateDone
)

var methods = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := &Request{
		RequestLine: RequestLine{
			HttpVersion:   "",
			RequestTarget: "",
			Method:        "",
		},
	}

	buf := make([]byte, 8)

	readToIndex := 0

	for request.state != requestStateDone {
		if readToIndex >= len(buf) {
			buf = append(buf, make([]byte, 8)...)
		}
		readSize, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		readToIndex += readSize
		parsedSize, err := request.parse(buf[:readToIndex])
		if parsedSize > 0 {
			copy(buf, buf[parsedSize:readToIndex])
			readToIndex -= parsedSize
		}
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == requestStateInitialized {
		requestLine, size, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if size == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.state = requestStateDone
		return size, nil
	}
	return 0, fmt.Errorf("state is not initilized.")
}

func parseRequestLine(data []byte) (RequestLine, int, error) {

	if bytes.Contains(data, []byte("\r\n")) {
		return RequestLine{}, 0, nil
	}

	parts := bytes.Split(data, []byte(" "))

	if len(parts) <= 2 {
		return RequestLine{}, 0, fmt.Errorf("request format is error.")
	}

	requestLine := RequestLine{
		Method:        "",
		RequestTarget: "",
		HttpVersion:   "",
	}

	actualMethod := string(parts[0])
	requestTarget := string(parts[1])
	httpVersion := string(parts[2])

	if !slices.Contains(methods, actualMethod) {
		return RequestLine{}, 0, fmt.Errorf("invalid method.")
	}

	requestLine.Method = actualMethod
	requestLine.RequestTarget = requestTarget

	versionParts := strings.Split(httpVersion, "/")

	if len(versionParts) != 2 {
		return RequestLine{}, 0, fmt.Errorf("http version is error.")
	}

	version, err := strconv.ParseFloat(versionParts[1], 64)

	if err != nil || version != 1.1 {
		return RequestLine{}, 0, fmt.Errorf("http version error.")
	}

	requestLine.HttpVersion = versionParts[1]

	return requestLine, len(data), nil
}
