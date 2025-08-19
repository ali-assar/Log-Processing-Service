package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ali-assar/Log-Processing-Service/mock-log-generator/internal/cli"
	wsserver "github.com/ali-assar/Log-Processing-Service/mock-log-generator/internal/ws"
)

func main() {
	log.Println("Starting Mock Log Generator...")
	// Parse CLI flags
	cfg, err := cli.Parse()
	if err != nil {
		log.Fatalf("failed to parse CLI args: %v", err)
	}

	mux := http.NewServeMux()
	wsserver.RegisterRoutes(mux, cfg)

	srv := &http.Server{
		Addr:              cfg.Url,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server in background
	go func() {
		log.Printf("Mock Log Generator listening on :%s",cfg.Url)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	log.Println("Shutdown signal received, stopping server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Mock Log Generator stopped")
}
