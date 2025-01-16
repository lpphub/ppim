package ext

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type BatchFunc[T any] func(t []T) error

type BatchProcessor[T any] struct {
	queue         chan T             // 用于接收数据的队列
	flushSignal   chan struct{}      // 用于触发批量处理的信号
	mu            sync.Mutex         // 用于保护批次的互斥锁
	wg            sync.WaitGroup     // 用于等待所有 goroutine 完成
	ctx           context.Context    // 上下文
	cancel        context.CancelFunc // 取消函数
	maxBatchSize  int                // 批次大小
	flushInterval time.Duration      // 刷新间隔
	workerCount   int                // 协程数量
	fn            BatchFunc[T]       // 批量处理函数
}

func NewBatchProcessor[T any](ctx context.Context, workerCount, maxBatchSize int, flushInterval time.Duration, batchFn BatchFunc[T]) *BatchProcessor[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &BatchProcessor[T]{
		ctx:           ctx,
		cancel:        cancel,
		queue:         make(chan T, maxBatchSize*workerCount),
		flushSignal:   make(chan struct{}),
		maxBatchSize:  maxBatchSize,
		flushInterval: flushInterval,
		workerCount:   workerCount,
		fn:            batchFn,
	}
}

func (p *BatchProcessor[T]) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			p.run(workerID)
		}(i)
	}
}

func (p *BatchProcessor[T]) Stop() {
	p.cancel() // 取消上下文
	close(p.queue)
	p.wg.Wait()
}

func (p *BatchProcessor[T]) run(_ int) {
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	var batch []T
	for {
		select {
		case t, ok := <-p.queue:
			if !ok {
				if len(batch) > 0 {
					p.flush(batch)
				}
				return
			}
			batch = append(batch, t)
			if len(batch) >= p.maxBatchSize {
				p.flush(batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				p.flush(batch)
				batch = nil
			}
		case <-p.flushSignal:
			if len(batch) > 0 {
				p.flush(batch)
				batch = nil
			}
		case <-p.ctx.Done():
			if len(batch) > 0 {
				p.flush(batch)
			}
			return
		}
	}
}

func (p *BatchProcessor[T]) flush(batch []T) {
	if err := p.fn(batch); err != nil {
		fmt.Printf("Batch processor flush error: %v\n", err)
		// 可以在这里添加重试逻辑
	}
}

func (p *BatchProcessor[T]) Submit(t T) error {
	select {
	case p.queue <- t:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("batch processor is stopped")
	}
}
