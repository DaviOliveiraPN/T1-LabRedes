package http1

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const idleTimeout = 5 * time.Second

type Server struct {
	Addr string
	Root string
}

func NewServer(port int, root string) *Server {
	return &Server{Addr: fmt.Sprintf(":%d", port), Root: root}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Printf("listening on %s, root %s", s.Addr, s.Root)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	var buf []byte

	for {
		conn.SetReadDeadline(time.Now().Add(idleTimeout))
		req, rest, err := readRequest(conn, buf)
		if errors.Is(err, ErrBadRequest) {
			res := errorResponse(400)
			res.Headers["Connection"] = "close"
			res.Write(conn, false)
			return
		}
		if err != nil {
			return
		}
		buf = rest

		log.Printf("%s %s %s", conn.RemoteAddr(), req.Method, req.Target)

		res := s.handle(req)
		keep := keepAlive(req)
		if keep {
			res.Headers["Connection"] = "keep-alive"
		} else {
			res.Headers["Connection"] = "close"
		}
		if err := res.Write(conn, req.Method == "HEAD"); err != nil || !keep {
			return
		}
	}
}

// HTTP/1.1 e persistente por padrao; HTTP/1.0 so com "keep-alive" explicito.
func keepAlive(req *Request) bool {
	c := strings.ToLower(req.Headers["connection"])
	if strings.Contains(c, "close") {
		return false
	}
	if req.Version == "HTTP/1.0" {
		return strings.Contains(c, "keep-alive")
	}
	return true
}
