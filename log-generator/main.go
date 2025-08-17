// mock_logs.go
package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Log struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Component string `json:"component"`
	TraceID   string `json:"trace_id"`
	SpanID    string `json:"span_id"`
	ParentID  string `json:"parent_id"`
}

var (
	levels = []string{
		"INFO", "WARN", "ERROR", "DEBUG", "TRACE",
		"FATAL", "CRITICAL", "PANIC", "UNKNOWN",
	}
	services   = []string{"auth", "orders", "payments", "notifications"}
	components = []string{
		"db", "api", "cache", "worker", "scheduler",
		"gateway", "service", "client", "middleware", "utils",
	}
	messages = []string{
		"user login successful", "user login failed",
		"order created", "payment declined",
		"cache miss", "db connection lost",
		"service started", "service stopped", "service restarted",
		"service crashed", "service recovered", "service updated",
		"request received", "request processed", "request failed",
		"request timed out", "request cancelled", "request started", "request completed",
		"system error", "system warning", "system info", "system debug",
		"system trace", "system fatal", "system panic", "system unknown",
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/logs/stream", func(w http.ResponseWriter, r *http.Request) {
		// Updated/extra headers for streaming
		w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")          // helps with nginx
		w.Header().Set("Access-Control-Allow-Origin", "*") // simple CORS for browser clients

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		// Optional fixed interval via ?interval_ms=NNN (10..10000)
		interval := time.Duration(rand.Intn(1000)) * time.Millisecond
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

		for {
			// Stop when the client goes away
			select {
			case <-ctx.Done():
				return
			default:
			}

			logEntry := newRandomLog()

			jsonData, err := json.Marshal(logEntry)
			if err != nil {
				continue
			}

			// Debug: mirror what's sent to the client
			log.Printf("sent log: %s", jsonData)

	
			_, _ = w.Write(jsonData)
			_, _ = w.Write([]byte("\n"))
			flusher.Flush()

			time.Sleep(interval)
		}
	})

	// Log server error instead of ignoring it
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newRandomLog() Log {
	return Log{
		Timestamp: time.Now().Unix(),
		Level:     pick(levels),
		Message:   pick(messages),
		Service:   pick(services),
		Component: pick(components),
		TraceID:   randomID(),
		SpanID:    randomID(),
		ParentID:  randomID(),
	}
}

func pick(options []string) string {
	return options[rand.Intn(len(options))]
}

func randomID() string {
	// 8 random bytes -> 16 hex chars; falls back to PRNG if needed
	b := make([]byte, 8)
	if _, err := crand.Read(b); err != nil {
		return fmt.Sprintf("%x", rand.Int63())
	}
	return hex.EncodeToString(b)
}

func printLog(entry Log) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("error marshaling log: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
