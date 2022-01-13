package main

import (
	"log"
	"net"
)

func main() {
	server := newServer()
	go server.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("OOOOOOOPS!")
	}
	defer listener.Close()

	log.Printf("Start!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Fatal, %s", err.Error())
			continue
		}
		go server.newClient(conn)
	}
}
