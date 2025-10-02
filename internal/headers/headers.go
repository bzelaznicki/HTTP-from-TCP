package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	i := bytes.Index(data, []byte(crlf))
	if i == -1 {
		return 0, false, nil
	}
	if i == 0 {
		return 2, true, nil
	}
	name, value, err := parseHeader(data[:i])
	if err != nil {
		return 0, false, err
	}
	h[name] = value
	return i + 2, false, nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	name := bytes.TrimLeft(parts[0], " \t")
	if len(name) == 0 || bytes.ContainsAny(name, " \t") {
		return "", "", fmt.Errorf("malformed field name")
	}
	if !validHeaderKey(name) {
		return "", "", fmt.Errorf("invalid header key")
	}
	name = bytes.ToLower(name)
	value := bytes.TrimSpace(parts[1])

	return string(name), string(value), nil
}

func validHeaderKey(b []byte) bool {
	for _, c := range b {
		switch {
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
		case c == '!' || c == '#' || c == '$' || c == '%' || c == '&' ||
			c == '\'' || c == '*' || c == '+' || c == '-' || c == '.' ||
			c == '^' || c == '_' || c == '`' || c == '|' || c == '~':
		default:
			return false
		}
	}
	return len(b) > 0
}
