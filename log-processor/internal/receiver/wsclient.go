package receiver

import (
	"context"
)

// Start connects to the generator WebSocket and streams messages.
// TODO: Implement WebSocket dial, read loop, decode JSON, and forward to parser/workerpool.
func Start(ctx context.Context, url string) error {
	// TODO: dial ws server (gorilla/websocket), loop on ReadMessage, decode, forward.
	// Return when ctx is done or on fatal error.
	<-ctx.Done()
	return ctx.Err()
}
