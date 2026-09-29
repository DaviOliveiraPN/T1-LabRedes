package http1

import (
	"errors"
	"os"
)

var ErrForbidden = errors.New("forbidden")

func (s *Server) handle(req *Request) *Response {
	if req.Method != "GET" && req.Method != "HEAD" {
		res := errorResponse(405)
		res.Headers["Allow"] = "GET, HEAD"
		return res
	}

	path, err := s.resolvePath(req.Path)
	if err != nil {
		return errorResponse(403)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		return errorResponse(404)
	}

	res := newResponse(200, body)
	res.Headers["Content-Type"] = contentType(path)
	return res
}

func (s *Server) resolvePath(p string) (string, error) {
	return "", ErrForbidden
}
