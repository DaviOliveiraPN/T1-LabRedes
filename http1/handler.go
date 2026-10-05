package http1

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

// resolvePath converte o caminho da URL num caminho dentro de Root.
// Qualquer ".." que suba acima da raiz vira ErrForbidden.
// A checagem e feita antes de limpar o caminho, senao "/../../x" viraria "/x" e passaria.
func (s *Server) resolvePath(p string) (string, error) {
	depth := 0
	for _, seg := range strings.Split(strings.ReplaceAll(p, "\\", "/"), "/") {
		switch seg {
		case "", ".":
		case "..":
			depth--
			if depth < 0 {
				return "", ErrForbidden
			}
		default:
			depth++
		}
	}

	full := filepath.Join(s.Root, filepath.FromSlash(p))
	if info, err := os.Stat(full); err == nil && info.IsDir() {
		full = filepath.Join(full, "index.html")
	}
	return full, nil
}
