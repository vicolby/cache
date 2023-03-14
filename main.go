package main

import "github.com/vicolby/cache/cache"

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	c := cache.New()
	server := NewServer(opts, c)
	server.Start()
}
