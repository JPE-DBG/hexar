package main

import (
	"flag"
	"hexar/internal/lobby"
	"hexar/internal/net"
	"log"
)

// Overridden at build time: go build -ldflags "-X main.version=<git-sha>"
var version = "dev"

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	clientDir := flag.String("client", "./client/dist", "path to client dist")
	flag.Parse()

	lob := lobby.New()
	srv := net.NewServer(lob, *clientDir, version)
	log.Fatal(net.ListenAndServe(*addr, srv))
}
