package main

import (
	"log"
	"socks555"
)

func main() {
	server := socks555.SOCKS5server{
		IP:   "localhost",
		Port: 1080,
	}
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
