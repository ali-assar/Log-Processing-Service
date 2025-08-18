package api

import "net/http"

// RegisterRoutes wires HTTP endpoints for stats and health checks.
func RegisterRoutes(mux *http.ServeMux) {
	// TODO: implement /healthz and /stats endpoints.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
