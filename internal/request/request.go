package request

import (
	"bytes"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
	State       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 1024

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	for !request.done() {
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}

		readToIndex += n
		readN, err := request.parse(buf[:readToIndex+n])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:readToIndex])
		readToIndex -= readN

	}
	return request, nil
}

func parseRequestLine(req []byte) (*RequestLine, int, error) {

	idx := bytes.Index(req, []byte(crlf))

	if idx == -1 {
		return nil, 0, nil
	}
	reqLine := req[:idx]
	read := idx + len(crlf)
	parts := bytes.Split(reqLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid number of values in request line")
	}
	method := parts[0]

	for _, c := range method {
		if c < 'A' || c > 'Z' {

			return nil, 0, fmt.Errorf("invalid method: %s", method)

		}
	}

	target := parts[1]

	verParts := bytes.Split(parts[2], []byte("/"))

	httpPart := verParts[0]

	if string(httpPart) != "HTTP" {
		return nil, 0, fmt.Errorf("unrecognized HTTP version: %s", httpPart)
	}

	ver := verParts[1]

	if string(ver) != "1.1" {
		return nil, 0, fmt.Errorf("unrecognized HTTP version: %s", ver)
	}
	requestData := RequestLine{
		Method:        string(method),
		RequestTarget: string(target),
		HttpVersion:   string(ver),
	}

	return &requestData, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.State = StateDone
		case StateDone:
			break outer
		}

	}
	return read, nil

}

func (r *Request) done() bool {
	return r.State == StateDone
}
