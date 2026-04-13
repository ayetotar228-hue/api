package worker

import (
	"api/pkg/broker/producer"
	"context"
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	workers  int
	jobQueue chan Job
	handlers map[string]JobHandler
	producer *producer.Producer
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewPool(workers int, producer *producer.Producer) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workers:  workers,
		jobQueue: make(chan Job, 100),
		handlers: make(map[string]JobHandler),
		producer: producer,
		ctx:      ctx,
		cancel:   cancel,
	}
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
		job.Result <- fmt.Errorf("no handler for type: %s", job.Type)
		return
	}

	err := handler(job.Context, job.Payload)

	job.Result <- err
}

func (p *Pool) Submit(job Job) error {
	if job.Result == nil {
		job.Result = make(chan error, 1)
	}
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("pool stopped")
	case p.jobQueue <- job:
		return <-job.Result
	default:
		return fmt.Errorf("queue full")
	}
}

func (p *Pool) Stop() {
	log.Println("Stopping worker pool...")
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
	log.Println("Worker pool stopped")
}
