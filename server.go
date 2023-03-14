package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

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
		go s.handleCommand(buf[:n], conn)

	}
}

func (s *Server) handleCommand(rawCmd []byte, conn net.Conn) {
	cmd, err := parseCommand(rawCmd)
	if err != nil {
		log.Printf("failed to parse command: %s \n", err)
		return
	}

	switch cmd.Command {
	case "SET":
		err := s.cache.Set(cmd.Key, cmd.Value, cmd.TTL)
		if err != nil {
			log.Printf("failed to set: %s \n", err)
			return
		}
		conn.Write([]byte("OK"))

	case "GET":
		val, err := s.cache.Get(cmd.Key)
		if err != nil {
			log.Printf("failed to set: %s \n", err)
			return
		}
		conn.Write(val)

	case "HAS":
		has := s.cache.Has(cmd.Key)
		if has {
			conn.Write([]byte("OK"))
		} else {
			conn.Write([]byte("NOT FOUND"))
		}
	case "DEL":
		err := s.cache.Delete(cmd.Key)
		if err != nil {
			log.Printf("failed to delete: %s \n", err)
			return
		}
		conn.Write([]byte("OK"))
	}
}

func parseCommand(rawCmd []byte) (*Message, error) {
	parts := strings.Split(string(rawCmd), " ")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid command")
	}

	cmd := parts[0]
	switch cmd {
	case "SET":
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid set command")
		}
		key := []byte(parts[1])
		value := []byte(parts[2])
		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid ttl")
		}
		return &Message{Command: Command(cmd), Key: key, Value: value, TTL: time.Duration(ttl) * time.Second}, nil
	case "GET":
		key := []byte(parts[1])
		return &Message{Command: Command(cmd), Key: key}, nil
	case "HAS":
		key := []byte(parts[1])
		return &Message{Command: Command(cmd), Key: key}, nil
	case "DEL":
		key := []byte(parts[1])
		return &Message{Command: Command(cmd), Key: key}, nil
	default:
		return nil, fmt.Errorf("invalid command %s", cmd)
	}
}
