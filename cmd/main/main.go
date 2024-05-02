package main

import (
	"log"
	"os"
	"strconv"
	"video-streaming/server"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	p, _ := strconv.Atoi(port)

	log.Fatal(server.StartServer(p, 45))

}
