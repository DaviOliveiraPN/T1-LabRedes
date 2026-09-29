package http1

import (
	"errors"
	"fmt"
	"log"
	"net"
)

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
		req, rest, err := readRequest(conn, buf)
		if errors.Is(err, ErrBadRequest) {
			errorResponse(400).Write(conn, false)
			return
		}
		if err != nil {
			return
		}
		buf = rest

		res := s.handle(req)
		if err := res.Write(conn, req.Method == "HEAD"); err != nil {
			return
		}

		if !keepAlive(req) {
			return
		}
	}
}

func keepAlive(req *Request) bool {
	return false
}
