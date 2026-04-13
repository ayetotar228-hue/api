package worker

import "context"

type Job struct {
	ID      int
	Type    string
	Payload interface{}
	Context context.Context
	Result  chan error
}

type JobHandler func(ctx context.Context, payload interface{}) error
