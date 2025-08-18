package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	api "github.com/ali-assar/Log-Processing-Service/log-processor/internal/api"
	receiver "github.com/ali-assar/Log-Processing-Service/log-processor/internal/receiver"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// TODO: make configurable via flags/env.
	generatorURL := "ws://localhost:8080/ws/logs"

	// Start the receiver (WebSocket client) in the background.
	go func() {
		if err := receiver.Start(ctx, generatorURL); err != nil {
			log.Printf("receiver stopped: %v", err)
		}
	}()

	// HTTP API (for future /stats, etc.) on a different port to avoid conflicts.
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:              ":8090",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("log-processor HTTP API listening on :8090")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}
