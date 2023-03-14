package main

import (
	"fmt"
	"log"
	"net"

	"github.com/vicolby/cache/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	opts  ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		opts:  opts,
		cache: c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.opts.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("Server listening on [%s] \n", s.opts.ListenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %s \n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("failed to read: %s \n", err)
			break
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}
