package main

import (
	"flag"
	"hexar/internal/net"
	"hexar/internal/room"
	"log"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	clientDir := flag.String("client", "./client/dist", "path to client dist")
	flag.Parse()

	r := room.New()
	go r.Run()

	srv := net.NewServer(r, *clientDir)
	log.Fatal(net.ListenAndServe(*addr, srv))
}
