package receiver

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ali-assar/Log-Processing-Service/log-processor/internal/parser"
	"github.com/ali-assar/Log-Processing-Service/log-processor/internal/workerpool"
	"github.com/gorilla/websocket"
)

// Start connects to the generator WebSocket and streams messages.
// TODO: Forward decoded messages to parser/workerpool.
func Start(ctx context.Context, url string, pool *workerpool.Pool) error {
	log.Printf("Connecting to WebSocket server at: %s", url)

	dialer := &websocket.Dialer{
		HandshakeTimeout:  5 * time.Second,
		EnableCompression: true,
	}

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Successfully connected to WebSocket server")

	// Close on context cancellation to unblock ReadMessage.
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
		log.Printf("Shutting down WebSocket connection...")
		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "shutdown"),
			time.Now().Add(1*time.Second),
		)
		_ = conn.Close()
	}()

	conn.SetPingHandler(func(data string) error {
		log.Printf("Received ping from server")
		return conn.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(data string) error {
		log.Printf("Received pong from server")
		return nil
	})

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	go func() {
		for {
			select {
			case <-pingTicker.C:
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
					log.Printf("Failed to send ping: %v", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	messageCount := 0
	var ignored, parseErrs, parsedOK, submitted, dropped int
	lastLog := time.Now()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			log.Printf("WebSocket read error: %v", err)
			return err
		}

		messageCount++

		var probe struct {
			Level string `json:"level"`
		}
		if err := json.Unmarshal(message, &probe); err == nil && probe.Level == "" {
			if messageCount <= 3 {
				log.Printf("ignoring non-log frame: %s", message)
			}
			ignored++
			// fall through to summary logging below
		} else {
			entry, err := parser.Parse(message)
			if err != nil {
				log.Printf("parse error: %v", err)
				parseErrs++
			} else {
				parsedOK++
				if pool.SubmitWithTimeout(entry, 10*time.Millisecond) {
					submitted++
				} else {
					dropped++
				}
			}
		}

		// Periodic summary: every 100 msgs or every 2s
		if messageCount%100 == 0 || time.Since(lastLog) > 2*time.Second {
			stats := pool.Stats()
			log.Printf("ws=%s total=%d ok=%d parse_err=%d ignored=%d submitted=%d dropped=%d queue=%d processed=%d workers=%d",
				url, messageCount, parsedOK, parseErrs, ignored, submitted, dropped, stats.Queue, stats.Processed, stats.Workers)
			lastLog = time.Now()
		}
	}
}
