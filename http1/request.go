package http1

import (
	"bytes"
	"errors"
	"net"
)

var ErrBadRequest = errors.New("bad request")

type Request struct {
	Method  string
	Target  string
	Path    string
	Version string
	Headers map[string]string
}

func readRequest(conn net.Conn, buf []byte) (*Request, []byte, error) {
	chunk := make([]byte, 4096)
	for {
		if i := bytes.Index(buf, []byte("\r\n\r\n")); i >= 0 {
			req, err := parseRequest(buf[:i])
			return req, buf[i+4:], err
		}
		n, err := conn.Read(chunk)
		buf = append(buf, chunk[:n]...)
		if err != nil {
			return nil, nil, err
		}
	}
}

func parseRequest(head []byte) (*Request, error) {
	return nil, ErrBadRequest
}

func decodePath(target string) (string, error) {
	return "", ErrBadRequest
}
