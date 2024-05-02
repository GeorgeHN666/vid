package main

import (
	"log"
	"video-streaming/server"
)

func main() {

	log.Fatal(server.StartServer(6354, 45))

}
