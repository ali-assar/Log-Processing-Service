package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	receiver "github.com/ali-assar/Log-Processing-Service/log-processor/internal/receiver"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	log.Println("Starting Log Processing Service...")
	// TODO: make configurable via flags/env.
	var generatorURL []string
	generatorURL = append(generatorURL, "ws://localhost:8080/ws/logs")

	receiverErr := make(chan error, 1)
	for _, v := range generatorURL {
		go func() {
			if err := receiver.Start(ctx, v); err != nil {
				log.Printf("Receiver stopped: %v", err)
				receiverErr <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received, waiting for graceful shutdown...")
		// Give some time for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		select {
		case <-shutdownCtx.Done():
			log.Println("Shutdown timeout reached")
		case err := <-receiverErr:
			log.Printf("Receiver stopped during shutdown: %v", err)
		}
	case err := <-receiverErr:
		log.Printf("Receiver failed: %v", err)
	}

	// HTTP API (for future /stats, etc.) on a different port to avoid conflicts.
	// mux := http.NewServeMux()
	// api.RegisterRoutes(mux)

	// srv := &http.Server{
	// 	Addr:              ":8080",
	// 	Handler:           mux,
	// 	ReadHeaderTimeout: 5 * time.Second,
	// }

	// log.Println("log-processor HTTP API listening on :8080")
	// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 	log.Fatalf("http server error: %v", err)
	// }
}
