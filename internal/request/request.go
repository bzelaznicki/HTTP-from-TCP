package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	parsedLine, err := parseRequestLine(req)

	if err != nil {
		return nil, err
	}
	request := Request{
		RequestLine: *parsedLine,
	}

	return &request, nil
}

func parseRequestLine(req []byte) (*RequestLine, error) {

	idx := bytes.Index(req, []byte(crlf))

	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	reqLine := string(req[:idx])

	parts := strings.Split(reqLine, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid number of values in request line")
	}
	method := parts[0]

	for _, c := range method {
		if c < 'A' || c > 'Z' {

			return nil, fmt.Errorf("invalid method: %s", method)

		}
	}

	target := parts[1]

	verParts := strings.Split(parts[2], "/")

	httpPart := verParts[0]

	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP version: %s", httpPart)
	}

	ver := verParts[1]

	if ver != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP version: %s", ver)
	}
	requestData := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   ver,
	}

	return &requestData, nil
}
