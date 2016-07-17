package main

import (
	"flag"
	"github.com/ilackarms/crawl/server"
	"log"
	"github.com/ilackarms/crawl/client"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9000", "address of server (client only)")
	name := flag.String("name", "", "name your character (must be unique)")
	serverMode := flag.Bool("server", false, "run in server mode")
	flag.Parse()
	if *serverMode {
		server.Start()
	} else {
		if *addr == "" {
			log.Fatal("must provide server address with -addr flag")
		}
		if *name == "" {
			log.Fatal("must provide name with -name flag")
		}
		client.Start(*name, *addr)
	}
}

