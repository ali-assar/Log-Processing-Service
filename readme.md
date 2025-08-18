## 🚀 Real-World Project: Log Processing Service

### 🎯 Project Goal

Build a **Log Processing Service** that simulates how real backend systems process logs concurrently, aggregate statistics, and serve results via an HTTP API.

This project demonstrates:

* Efficient concurrency with worker pools.
* Real-world patterns (fan-in, fan-out, cancellations).
* Integration with a database.
* Exposing results via REST API.

---

## 🏗️ System Design

### **Architecture**

1. **Input Layer**

   * Logs are fed into the system (from files or API).

2. **Processing Layer (Worker Pool)**

   * Multiple workers parse log entries concurrently.
   * Each worker extracts error counts, warnings, etc.

3. **Aggregation Layer (Fan-In)**

   * Collect results from workers.
   * Update database with aggregated statistics.

4. **Query Layer (HTTP Server)**

   * REST endpoints to fetch statistics (e.g., total errors, error frequency, logs per service).

---

## 📋 Features

1. **Concurrent Log Processing**

   * Goroutines + channels for job distribution.
   * Worker pool to limit concurrency.

2. **Error Counting & Storage**

   * Parse logs for error/warning/info levels.
   * Store results in SQLite/Postgres (start simple with SQLite).

3. **Graceful Shutdown**

   * Use `context.WithCancel` to stop workers.
   * Handle OS signals (`SIGINT`, `SIGTERM`).

4. **HTTP Query API**

   * Endpoint: `/stats` → Returns aggregated error counts.
   * Endpoint: `/stats/:service` → Service-specific stats.

---

## 📂 Project Structure

```
log-processing-service/
│── cmd/
│   └── main.go           # Entry point
│
│── internal/
│   ├── workerpool/       # Worker pool implementation
│   ├── parser/           # Log parsing logic
│   ├── storage/          # Database layer
│   └── api/              # HTTP server & handlers
│
│── pkg/
│   └── models/           # Data models (LogEntry, Stats)
│
│── configs/
│   └── config.yaml       # Configs (DB, server port, etc.)
│
│── logs/                 # Sample input log files
│── go.mod
│── README.md
```

---

## 🛠️ Tech Stack

* **Go** (Concurrency, net/http, context)
* **SQLite/Postgres** (Persistent storage)
* **Docker** (Optional: for containerized deployment)
* **Makefile** (Optional: build/run automation)

---

## 📅 Suggested Weekly Breakdown

### **Week 1: Concurrency Foundations**

* Refresh goroutines, channels, select.
* Implement small exercises (fan-in, fan-out, worker pools).

### **Week 2: Worker Pool & Log Parsing**

* Implement worker pool.
* Write parser for log files (simple regex or structured).
* Test parsing with multiple workers.

### **Week 3: Aggregation & Storage**

* Add fan-in aggregator.
* Connect workers to SQLite.
* Store error counts and statistics.

### **Week 4: HTTP API & Graceful Shutdown**

* Expose `/stats` endpoint.
* Add service-level queries.
* Implement context cancellations for shutdown.
* Final integration + tests.

---

## ✅ Deliverables

* A working **log processing service**.
* Ability to process multiple logs concurrently.
* Aggregated stats stored in a DB.
* REST API to query results.
* Proper handling of cancellations and shutdowns.

---

## 🧪 Local Mock Log Generator

A lightweight WebSocket stream to feed the system during development.

- Endpoint: WS ws://localhost:8080/ws/logs
- Message format: JSON per message (one log entry)
- Query params:
  - interval_ms: integer (10..10000). Default: random up to 1000ms.
  - service: fixed service name override (e.g., auth).
  - level: fixed level override (INFO|WARN|ERROR).
  - level_weights: weighted distribution for levels. Format: INFO:70,WARN:20,ERROR:10
  - debug: if present, logs each emitted line on the server.

Examples:
- websocat ws://localhost:8080/ws/logs?interval_ms=200&level_weights=INFO:80,WARN:15,ERROR:5
- websocat ws://localhost:8080/ws/logs?service=payments&level=ERROR&debug=1

Sample message:
{"timestamp": 1735668123456, "level": "INFO", "message": "order created", "service": "orders", "component": "api", "trace_id": "...", "span_id": "...", "parent_id": "..."}