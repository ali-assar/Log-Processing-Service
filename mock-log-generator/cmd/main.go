package main

import (
	"log"
	"net/http"

	wsserver "github.com/ali-assar/Log-Processing-Service/mock-log-generator/internal/ws"
)

func main() {
	mux := http.NewServeMux()
	wsserver.RegisterRoutes(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
