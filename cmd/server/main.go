package main

import (
	"context"
	"go_redis/internal/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	srv := server.NewCacheServer(":3001")
	log.Println("Listening to server on port 3001...")
	wg.Add(1)

	go srv.Start(&wg)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	<-ctx.Done()

	srv.Quit()
	wg.Wait()
	os.Exit(0)
}
