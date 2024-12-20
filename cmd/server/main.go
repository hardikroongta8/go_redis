package main

import (
	"go_redis/internal/server"
	"log"
)

func main() {
	srv := server.NewCacheServer(":3001")
	log.Println("Listening to server on port 3001...")
	log.Fatal(srv.Start())

}
