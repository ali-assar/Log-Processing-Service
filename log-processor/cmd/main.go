package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ali-assar/Log-Processing-Service/log-processor/internal/cli"
	receiver "github.com/ali-assar/Log-Processing-Service/log-processor/internal/receiver"
	"github.com/ali-assar/Log-Processing-Service/log-processor/internal/workerpool"
)

type inMemoryStore struct {
	mu     sync.Mutex
	counts map[string]map[string]int
}

func (s *inMemoryStore) IncrementLevelCount(svc, lvl string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.counts[svc] == nil {
		s.counts[svc] = map[string]int{}
	}
	s.counts[svc][lvl]++
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create storage (temporary in-memory while learning)
	store := &inMemoryStore{counts: map[string]map[string]int{}}

	// Start pool: e.g., 8 workers, queue size 1024
	pool := workerpool.Start(ctx, 50, 2048, store)
	defer pool.Close()

	cfg, err := cli.Parse()
	if err != nil {
		log.Fatalf("failed to parse CLI args: %v", err)
	}

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s := pool.Stats()
				log.Printf("POOL stats: processed=%d queue=%d workers=%d", s.Processed, s.Queue, s.Workers)
			}
		}
	}()

	errCh := make(chan error, len(cfg.Urls))
	for _, u := range cfg.Urls {
		u := u
		go func() {
			if err := receiver.Start(ctx, u, pool); err != nil {
				errCh <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
	case err := <-errCh:
		log.Printf("receiver error: %v", err)
	}
	// pool.Close() deferred
	_ = time.Second
}
