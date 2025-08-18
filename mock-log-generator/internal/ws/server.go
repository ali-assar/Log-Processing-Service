package ws

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ali-assar/Log-Processing-Service/mock-log-generator/internal/generator"
	"github.com/ali-assar/Log-Processing-Service/mock-log-generator/internal/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Simple dev-only CORS; adjust for production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws/logs", WSLogsHandler)
}

func WSLogsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// interval_ms param (10..10000); default: random up to 1000ms
	interval := time.Duration(rng.Intn(1000)) * time.Millisecond
	if s := r.URL.Query().Get("interval_ms"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			if v < 10 {
				v = 10
			}
			if v > 10000 {
				v = 10000
			}
			interval = time.Duration(v) * time.Millisecond
		}
	}

	fixedService := strings.TrimSpace(r.URL.Query().Get("service"))
	fixedLevel := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("level")))
	debug := hasQueryKey(r, "debug")

	weights := generator.ParseLevelWeights(r.URL.Query().Get("level_weights"), types.DefaultLevelWeights)
	levelPicker := generator.NewWeightedPicker(weights, rng)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			entry := generator.NewRandomLog(rng, func() string {
				if fixedLevel != "" {
					return fixedLevel
				}
				return levelPicker.Pick()
			}, func() string {
				if fixedService != "" {
					return fixedService
				}
				return generator.PickFrom(rng, types.Services)
			})

			b, err := json.Marshal(entry)
			if err != nil {
				continue
			}
			if debug {
				log.Printf("sent log: %s", b)
			}
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				return // client disconnected or error
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func hasQueryKey(r *http.Request, key string) bool {
	_, ok := r.URL.Query()[key]
	return ok
}
