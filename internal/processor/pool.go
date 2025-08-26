package processor

import (
	"context"
	"sync"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
)

type Pool struct {
	Orders    chan models.Order
	Results   chan models.Order
	Wg        sync.WaitGroup
	Ctx       context.Context
	Cancel    context.CancelFunc
	Processed int
	Workers   int
}

func Start(ctx context.Context, workers, buf int) *Pool {
	ctx, cancel := context.WithCancel(ctx)
	pool := &Pool{
		Orders:  make(chan models.Order, buf),
		Results: make(chan models.Order, buf),
		Workers: workers,
		Ctx:     ctx,
		Cancel:  cancel,
	}

	for i := 0; i < workers; i++ {
		pool.Wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func Close(pool *Pool) {
	pool.Cancel()
	pool.Wg.Wait()
	close(pool.Orders)
	close(pool.Results)
}

func (p *Pool) worker(id int) {
	defer p.Wg.Done()
	for {
		select {
		case <-p.Ctx.Done():
			return
		case order, ok := <-p.Orders:
			if !ok {
				return
			}
			// implement logic before sending the result
			p.Results <- order
			p.Processed++
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
		processed: p.Processed,
		queue:     len(p.Orders),
		workers:   p.Workers,
	}
}
