package receiver

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ali-assar/Log-Processing-Service/log-processor/pkg/models"
	"github.com/gorilla/websocket"
)

// Start connects to the generator WebSocket and streams messages.
// TODO: Forward decoded messages to parser/workerpool.
func Start(ctx context.Context, url string) error {
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

		var logEntry models.LogEntry
		if err := json.Unmarshal(message, &logEntry); err != nil {
			var connMsg map[string]interface{}
			if err := json.Unmarshal(message, &connMsg); err == nil {
				if status, ok := connMsg["status"].(string); ok && status == "connected" {
					log.Printf("Connection confirmed. Server interval: %v ms", connMsg["interval_ms"])
					continue
				}
			}

			// If neither, log as raw message
			log.Printf("Received raw message: %s", message)
			continue
		}

		// Process the log entry
		log.Printf("[%s] %s/%s: %s", logEntry.Timestamp, logEntry.Service, logEntry.Level, logEntry.Message)

		// TODO: Forward to parser/workerpool for further processing
		// For now, just log the received entries
		if messageCount%100 == 0 {
			log.Printf("Processed %d messages so far", messageCount)
		}
	}
}
