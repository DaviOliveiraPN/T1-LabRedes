package http1

import (
	"bytes"
	"errors"
	"net"
	"net/url"
	"strings"
)

const maxHeaderBytes = 64 * 1024

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
		if len(buf) > maxHeaderBytes {
			return nil, nil, ErrBadRequest
		}
		n, err := conn.Read(chunk)
		buf = append(buf, chunk[:n]...)
		if err != nil {
			return nil, nil, err
		}
	}
}

func parseRequest(head []byte) (*Request, error) {
	lines := strings.Split(string(head), "\r\n")

	parts := strings.Split(lines[0], " ")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, ErrBadRequest
	}
	if parts[2] != "HTTP/1.1" && parts[2] != "HTTP/1.0" {
		return nil, ErrBadRequest
	}

	path, err := decodePath(parts[1])
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method:  parts[0],
		Target:  parts[1],
		Path:    path,
		Version: parts[2],
		Headers: map[string]string{},
	}

	for _, line := range lines[1:] {
		name, value, ok := strings.Cut(line, ":")
		if !ok || strings.TrimSpace(name) == "" {
			return nil, ErrBadRequest
		}
		req.Headers[strings.ToLower(strings.TrimSpace(name))] = strings.TrimSpace(value)
	}
	return req, nil
}

// decodePath tira a query string e decodifica o percent-encoding (%20 -> espaco).
func decodePath(target string) (string, error) {
	raw, _, _ := strings.Cut(target, "?")
	if !strings.HasPrefix(raw, "/") {
		return "", ErrBadRequest
	}
	p, err := url.PathUnescape(raw)
	if err != nil || strings.ContainsRune(p, 0) {
		return "", ErrBadRequest
	}
	return p, nil
}
