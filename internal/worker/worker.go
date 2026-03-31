package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	workers  int
	jobQueue chan Job
	handlers map[string]JobHandler
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewPool(workers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		workers:  workers,
		jobQueue: make(chan Job, 100),
		handlers: make(map[string]JobHandler),
		ctx:      ctx,
		cancel:   cancel,
	}

	return p
}

func (p *Pool) Register(jobType string, handler JobHandler) {
	p.handlers[jobType] = handler
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("Worker pool started with %d workers", p.workers)
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("Worker %d stopped", id)
			return
		case job, ok := <-p.jobQueue:
			if !ok {
				return
			}
			p.processJob(id, job)
		}
	}
}

func (p *Pool) processJob(workerID int, job Job) {
	handler, exists := p.handlers[job.Type]
	if !exists {
		log.Printf("Worker %d: no handler for job type %s", workerID, job.Type)
		return
	}

	if err := handler(job.Context, job); err != nil {
		log.Printf("Worker %d: job %d (%s) failed: %v", workerID, job.ID, job.Type, err)
		return
	}

	log.Printf("Worker %d: job %d (%s) completed", workerID, job.ID, job.Type)
}

func (p *Pool) Submit(job Job) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is stopped")
	default:
		select {
		case p.jobQueue <- job:
			return nil
		default:
			return fmt.Errorf("job queue is full")
		}
	}
}

func (p *Pool) Stop() {
	log.Println("Stopping worker pool...")
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
	log.Println("Worker pool stopped")
}
