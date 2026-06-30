package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	parseableHeaderBytes, _, ok := bytes.Cut(data, []byte(CRLF))

	if !ok {
		return 0, false, nil
	}

	if len(parseableHeaderBytes) == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(parseableHeaderBytes, []byte(":"), 2)

	if len(parts) != 2 {
		return 0, false, fmt.Errorf("header is not correct format")
	}

	key := string(parts[0])

	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("malformed header name")
	}

	if !isValidFieldName(key) {
		return 0, false, fmt.Errorf("Error field name")
	}

	actualValue := strings.TrimSpace(string(parts[1]))

	actualHeader := strings.ToLower(key)

	_, ok = h[actualHeader]

	if ok {
		h[actualHeader] = fmt.Sprintf("%v, %v", h[actualHeader], actualValue)
		return len(parseableHeaderBytes) + 2, false, nil
	}

	h[actualHeader] = actualValue

	return len(parseableHeaderBytes) + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	value, ok := h[strings.ToLower(key)]
	return value, ok
}

func (h Headers) Override(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Delete(key string) {
	delete(h, key)
}

func isValidFieldName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, c := range []byte(name) {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z':
		return true
	case c >= 'a' && c <= 'z':
		return true
	case c >= '0' && c <= '9':
		return true
	}
	switch c {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false
}
