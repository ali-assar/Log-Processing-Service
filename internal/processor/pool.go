package processor

import (
	"context"
	"sync"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
)

type Pool struct {
	orders    chan models.Order
	results   chan models.Order
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	processed int
	workers   int
}

func Start(ctx context.Context, workers, buf int) *Pool {
	ctx, cancel := context.WithCancel(ctx)
	pool := &Pool{
		orders:  make(chan models.Order, buf),
		results: make(chan models.Order, buf),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}

	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func Close(pool *Pool) {
	pool.cancel()
	pool.wg.Wait()
	close(pool.orders)
	close(pool.results)
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case order, ok := <-p.orders:
			if !ok {
				return
			}
			// implement logic before sending the result
			p.results <- order
			p.processed++
		}
	}
}

type status struct {
	processed int
	queue     int
	workers   int
}

func (p *Pool) Stats() status {
	return status{
		processed: p.processed,
		queue:     len(p.orders),
		workers:   p.workers,
	}
}
