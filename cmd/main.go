package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/handler"
	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/processor"
)

func main() {

	mux := http.NewServeMux()
	pool := processor.Start(context.Background(), 10, 100)
	defer processor.Close(pool)

	go func() {
		for result := range pool.Results {
			log.Printf("Order %s processed", result.ID)
		}
	}()

	handler.RegisterRoutes(mux, pool)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("API listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
