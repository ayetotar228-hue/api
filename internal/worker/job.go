package worker

import "context"

type Job struct {
	ID      int
	Type    string
	Payload interface{}
	Context context.Context
}

type JobHandler func(ctx context.Context, job Job) error
