package types

type Log struct {
	Timestamp int64  `json:"timestamp"` // Unix milliseconds
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Component string `json:"component"`
	TraceID   string `json:"trace_id"`
	SpanID    string `json:"span_id"`
	ParentID  string `json:"parent_id"`
}

var (
	// Reduced to the core levels your processor aggregates.
	Levels = []string{"INFO", "WARN", "ERROR"}

	// Default weights for level distribution (overridable via query).
	DefaultLevelWeights = map[string]int{
		"INFO":  70,
		"WARN":  20,
		"ERROR": 10,
	}

	Services   = []string{"auth", "orders", "payments", "notifications"}
	Components = []string{
		"db", "api", "cache", "worker", "scheduler",
		"gateway", "service", "client", "middleware", "utils",
	}
	Messages = []string{
		"user login successful", "user login failed",
		"order created", "payment declined",
		"cache miss", "db connection lost",
		"service started", "service stopped", "service restarted",
		"service crashed", "service recovered", "service updated",
		"request received", "request processed", "request failed",
		"request timed out", "request cancelled", "request started", "request completed",
		"system error", "system warning", "system info",
	}
)
