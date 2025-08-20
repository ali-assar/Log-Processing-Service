package api

import (
	"encoding/json"
	"net/http"

	"github.com/ali-assar/Log-Processing-Service/log-processor/internal/workerpool"
)

func RegisterRoutes(mux *http.ServeMux, statsFn func() workerpool.Stats) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		s := statsFn()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(s)
	})
}
