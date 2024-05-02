package main

import (
	"log"
	"video-streaming/server"
)

func main() {

	log.Fatal(server.StartServer(3000, 45))

}
