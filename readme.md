# Log Processing Service

A learning-focused, concurrent log processing pipeline composed of:
- mock-log-generator: a WebSocket server that emits JSON log messages.
- log-processor: a client that connects to one or more WebSocket sources, parses messages, and processes them via a bounded worker pool. It exposes a tiny HTTP API for health and stats.

## Architecture

WebSocket Sources (fan-in) -> Receiver (wsclient) -> Parser (validation/normalize) -> Worker Pool (bounded, N workers) -> Storage (in-memory; pluggable for SQLite)

Key properties:
- Fan-in: multiple websocket URLs share one pool.
- Backpressure: bounded queue; overload can drop with timeout.
- Graceful shutdown: context cancellation closes sockets and drains workers.
- Health/Stats: HTTP server publishes pool stats.

## Requirements

- Go 1.21+ (recommended 1.22+)
- Windows PowerShell examples below; adapt paths for your OS.

## Quick Start

Terminal 1: start the mock generator
- Default: listen on :8080 and emit logs at a fixed or random interval.

```powershell
go run .\mock-log-generator\cmd\main.go --url :8080 --interval-ms 250
# Per-connection override:
# ws://localhost:8080/ws/logs?interval_ms=50
```

Terminal 2: start the processor with one or more sources
```powershell
go run .\log-processor\cmd\main.go --urls "ws://localhost:8080/ws/logs" --http-addr ":9090"
# Multiple sources:
# go run .\log-processor\cmd\main.go --urls "ws://localhost:8080/ws/logs,ws://localhost:9090/ws/logs" --http-addr ":9090"
```

Check health and stats:
```powershell
curl http://localhost:9090/healthz
curl http://localhost:9090/stats
```

Example /stats response:
```json
{
  "Processed": 12345,
  "Queue": 12,
  "Workers": 50
}
```

## CLI Reference

mock-log-generator
- --url, -u: HTTP listen address (e.g., :8080)
- --interval-ms, -i: default interval in ms (0 = random per connection; clients may override via ?interval_ms=NN)

log-processor
- --urls, -u: comma-separated WebSocket URLs (e.g., ws://localhost:8080/ws/logs,ws://localhost:9090/ws/logs)
- --http-addr, -a: HTTP listen address for the stats API (default :9090)

## Runtime Logs

Receiver connection
- “Connecting to WebSocket server at: …”
- “Successfully connected to WebSocket server”
- Ping/pong keepalive logged at debug level.

Per-connection summaries (every 100 msgs or ~2s)
- “ws=<url> total=<N> ok=<parsedOK> parse_err=<parseErrs> ignored=<ignored> submitted=<submitted> dropped=<dropped> queue=<q> processed=<p> workers=<w>”

Notes
- The generator sends an initial handshake JSON without “level”. The receiver probes payloads and “ignored” counts reflect such control frames.
- Parser validates level (INFO/WARN/ERROR), timestamp sanity (+/- 7d), and non-empty service/message.

## Concurrency & Backpressure

Worker Pool
- N workers consume from a buffered jobs channel.
- pool.SubmitWithTimeout(entry, 10*time.Millisecond) applies backpressure: under load, submissions may be dropped to keep the system responsive.

Tune
- Increase workers to raise throughput until CPU or storage becomes the bottleneck.
- Adjust queue size to absorb bursts.
- The in-memory storage uses a mutex to safely aggregate counts.

Shutdown
- Ctrl+C triggers context cancellation:
  - Receivers close WebSocket connections (normal close frame).
  - Pool cancels workers and waits for them to finish.

## Project Layout

- mock-log-generator/
  - internal/cli: flags for listen address and default interval
  - internal/ws: WebSocket handler that emits logs (with optional per-connection interval)
  - cmd/main.go: server bootstrap
- log-processor/
  - internal/cli: URLs and API listen address
  - internal/receiver: wsclient.go (connect, read, probe, parse, submit)
  - internal/parser: parser.go (JSON -> models.LogEntry with validation)
  - internal/workerpool: pool.go (bounded queue, N workers, stats)
  - internal/api: server.go (/healthz, /stats)
  - cmd/main.go: wiring (pool, receivers, HTTP API)

Models
- pkg/models: LogEntry struct used across the pipeline.

## Troubleshooting

- “parse error: invalid level”: often the initial handshake frame (no “level”). The receiver ignores such frames; if you still see errors, inspect payload via temporary logs to confirm field names.
- Drops increasing: raise queue size, worker count, or relax Submit timeout. If storage becomes a hotspot, consider batching in workers or swapping to a DB-backed storage.
- “concurrent map writes”: ensure storage uses a mutex (the in-memory example already does).

## Next Steps

- Replace in-memory storage with SQLite (implement the Storage interface).
- Add batching/flush intervals if moving to a DB.
- Expose aggregated counters via the HTTP API.
- Add unit tests for parser and pool behavior under load.
