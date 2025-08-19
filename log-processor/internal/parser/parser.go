package parser

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ali-assar/Log-Processing-Service/log-processor/pkg/models"
)

// TODO: Define Parse([]byte) -> (LogEntry, error) and helpers.
// This layer should validate and normalize logs before aggregation.

var allowed = map[string]struct{}{"INFO": {}, "WARN": {}, "ERROR": {}}

func Parse(raw []byte) (models.LogEntry, error) {
	var e models.LogEntry

	if err := json.Unmarshal(raw, &e); err != nil {
		return models.LogEntry{}, err
	}

	e.Level = strings.ToUpper(strings.TrimSpace(e.Level))
	if _, ok := allowed[e.Level]; !ok {
		return models.LogEntry{}, errors.New("invalid level")
	}

	now := time.Now().UnixMilli()
	week := int64(7 * 24 * time.Hour / time.Millisecond)
	if e.Timestamp < now-week || e.Timestamp > now+week {
		return models.LogEntry{}, errors.New("timestamp out of range")
	}

	if strings.TrimSpace(e.Service) == "" || strings.TrimSpace(e.Message) == "" {
		return models.LogEntry{}, errors.New("missing service/message")
	}

	return e, nil
}
