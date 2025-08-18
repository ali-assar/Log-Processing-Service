package models

// LogEntry mirrors the generator message format.
type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Component string `json:"component"`
	TraceID   string `json:"trace_id"`
	SpanID    string `json:"span_id"`
	ParentID  string `json:"parent_id"`
}

// Stats is a placeholder for aggregated results.
type Stats struct {
	Total   int            `json:"total"`
	ByLevel map[string]int `json:"by_level"`
	BySvc   map[string]int `json:"by_service"`
}
