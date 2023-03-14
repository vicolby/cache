package main

import (
	"log"
	"net"
	"time"

	"github.com/vicolby/cache/cache"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(2 * time.Second)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatalf("failed to connect to server: %v", err)
		}
		conn.Write([]byte("GET Bar Foo 1000"))
	}()

	c := cache.New()
	server := NewServer(opts, c)
	server.Start()
}
