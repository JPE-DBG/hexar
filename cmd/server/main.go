package main

import (
	"flag"
	"hexar/internal/lobby"
	"hexar/internal/net"
	"log"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	clientDir := flag.String("client", "./client/dist", "path to client dist")
	flag.Parse()

	lob := lobby.New()
	srv := net.NewServer(lob, *clientDir)
	log.Fatal(net.ListenAndServe(*addr, srv))
}
