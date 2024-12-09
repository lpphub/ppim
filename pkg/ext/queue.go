package ext

import (
	"context"
	"sync"
)

type Task func(context.Context) error

type Queue struct {
	tasks   chan Task
	workers int
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func New(workers int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		tasks:   make(chan Task, 100), // buffer size of 100
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (q *Queue) Start() {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

func (q *Queue) Stop() {
	q.cancel()
	close(q.tasks)
	q.wg.Wait()
}

func (q *Queue) Add(task Task) {
	select {
	case <-q.ctx.Done():
		return
	case q.tasks <- task:
	}
}

// worker processes tasks from the queue
func (q *Queue) worker() {
	defer q.wg.Done()

	for {
		select {
		case <-q.ctx.Done():
			return
		case task, ok := <-q.tasks:
			if !ok {
				return
			}
			_ = task(q.ctx)
		}
	}
}
