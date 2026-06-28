package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"http.xlkv.io/internal/headers"
)

type parserState int

const (
	requestStateInitialized parserState = iota
	requestStateDone
	requestStateParsingHeaders
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
	Headers     headers.Headers
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
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
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
				if request.state != requestStateDone {
					return nil, fmt.Errorf("reqeust format error")
				}
				break
			}
			return nil, err
		}
		readToIndex += readSize
		parsedSize, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if parsedSize > 0 {
			copy(buf, buf[parsedSize:readToIndex])
			readToIndex -= parsedSize
		}
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	countParsedBytes := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[countParsedBytes:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		countParsedBytes += n
	}
	return countParsedBytes, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, size, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if size == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.state = requestStateParsingHeaders
		return size, nil
	case requestStateParsingHeaders:
		size, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if size == 0 {
			return 0, nil
		}
		if done {
			r.state = requestStateDone
		}
		return size, nil
	}
	return 0, fmt.Errorf("state is not initilized.")
}

func parseRequestLine(data []byte) (RequestLine, int, error) {

	idx := bytes.Index(data, []byte("\r\n"))

	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	parts := bytes.Split(data[:idx], []byte(" "))

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

	return requestLine, len(data[:idx+2]), nil
}
