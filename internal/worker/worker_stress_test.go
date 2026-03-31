package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_StressManyJobs(t *testing.T) {
	pool := NewPool(10)
	pool.Start()
	defer pool.Stop()

	var processed int32
	const totalJobs = 500

	pool.Register("stress_job", func(ctx context.Context, job Job) error {
		atomic.AddInt32(&processed, 1)
		time.Sleep(time.Microsecond * 10)
		return nil
	})

	for i := 0; i < totalJobs; i++ {
		err := pool.Submit(Job{ID: i, Type: "stress_job"})
		if err != nil {
			time.Sleep(time.Millisecond)
			err = pool.Submit(Job{ID: i, Type: "stress_job"})
		}
		assert.NoError(t, err)
	}

	time.Sleep(time.Second * 2)
	assert.Equal(t, int32(totalJobs), atomic.LoadInt32(&processed))
}
func TestWorkerPool_StressConcurrentSubmit(t *testing.T) {
	pool := NewPool(5)
	pool.Start()
	defer pool.Stop()

	var processed int32
	const totalJobs = 100

	pool.Register("concurrent_job", func(ctx context.Context, job Job) error {
		atomic.AddInt32(&processed, 1)
		return nil
	})

	var wg sync.WaitGroup
	for g := 0; g < 5; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < totalJobs/5; i++ {
				pool.Submit(Job{ID: i, Type: "concurrent_job"})
				time.Sleep(time.Microsecond * 100)
			}
		}()
	}

	wg.Wait()
	time.Sleep(time.Second * 2)

	assert.Equal(t, int32(totalJobs), atomic.LoadInt32(&processed))
}
func TestWorkerPool_StressSlowHandler(t *testing.T) {
	pool := NewPool(3)
	pool.Start()
	defer pool.Stop()

	var processed int32

	pool.Register("slow_job", func(ctx context.Context, job Job) error {
		time.Sleep(time.Millisecond * 50)
		atomic.AddInt32(&processed, 1)
		return nil
	})

	for i := 0; i < 100; i++ {
		pool.Submit(Job{ID: i, Type: "slow_job"})
	}

	time.Sleep(time.Second * 3)
	assert.Equal(t, int32(100), atomic.LoadInt32(&processed))
}

func TestWorkerPool_StressQueueFull(t *testing.T) {
	pool := NewPool(1)
	pool.Start()
	defer pool.Stop()

	pool.Register("blocking_job", func(ctx context.Context, job Job) error {
		time.Sleep(time.Second)
		return nil
	})

	var errors int
	for i := 0; i < 200; i++ {
		err := pool.Submit(Job{ID: i, Type: "blocking_job"})
		if err != nil {
			errors++
		}
	}

	assert.Greater(t, errors, 0, "Some jobs should be rejected when queue is full")
}
