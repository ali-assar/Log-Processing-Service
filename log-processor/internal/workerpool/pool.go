package workerpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ali-assar/Log-Processing-Service/log-processor/pkg/models"
)

// TODO: Implement a bounded worker pool that accepts parsed LogEntry jobs.
// Consider Start(ctx, n int), Submit(job), and Close() patterns.

// interface to plug sqlite
type Storage interface {
	IncrementLevelCount(service, level string) error
}

type Pool struct {
	jobs      chan models.LogEntry
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	store     Storage
	processed int64
	workers   int
}

func Start(ctx context.Context, workers, buf int, store Storage) *Pool {
	cctx, cancel := context.WithCancel(ctx)
	p := &Pool{
		jobs:    make(chan models.LogEntry, buf),
		ctx:     cctx,
		cancel:  cancel,
		store:   store,
		workers: workers,
	}

	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	return p
}

func (p *Pool) SubmitWithTimeout(e models.LogEntry, d time.Duration) bool {
	select {
	case <-p.ctx.Done():
		return false
	case p.jobs <- e:
		return true
	case <-time.After(d):
		return false
	}
}

func (p *Pool) Close() {
	p.cancel()
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case e, ok := <-p.jobs:
			if !ok {
				return
			}
			_ = p.store.IncrementLevelCount(e.Service, e.Level)
			atomic.AddInt64(&p.processed, 1)
		}
	}
}

type Stats struct {
	Processed int64
	Queue     int
	Workers   int
}

func (p *Pool) Stats() Stats {
	return Stats{
		Processed: atomic.LoadInt64(&p.processed),
		Queue:     len(p.jobs),
		Workers:   p.workers,
	}
}
