package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assignmentap/internal/server"
	"assignmentap/internal/store"
	"assignmentap/internal/worker"
)

func main() {
	st := store.NewStore[string, string]()
	srv := server.NewServer(st)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	stopWorker := make(chan struct{})
	worker.StartWorker(stopWorker, srv)

	go func() {
		log.Println("Server running on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	close(stopWorker)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)

	log.Println("Server stopped")
}
