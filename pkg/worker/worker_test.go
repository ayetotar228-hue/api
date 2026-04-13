package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_SubmitAndProcess(t *testing.T) {
	pool := NewPool(2)
	pool.Start()
	defer pool.Stop()

	var processed int32

	pool.Register("test_job", func(ctx context.Context, job Job) error {
		atomic.AddInt32(&processed, 1)
		return nil
	})

	for i := 0; i < 5; i++ {
		err := pool.Submit(Job{ID: i, Type: "test_job"})
		assert.NoError(t, err)
	}

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(5), atomic.LoadInt32(&processed))
}

func TestWorkerPool_Stop(t *testing.T) {
	pool := NewPool(2)
	pool.Start()

	pool.Stop()

	err := pool.Submit(Job{ID: 1, Type: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stopped")
}
