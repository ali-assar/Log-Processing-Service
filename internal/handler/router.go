package handler

import (
	"net/http"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/processor"
)

func RegisterRoutes(router *http.ServeMux, pool *processor.Pool) {
	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		CreateOrderHandler(w, r, pool)
	})
}
