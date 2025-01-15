package ext

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	BatchDefaultMaxSize       = 100              // 默认内存队列的最大容量
	BatchDefaultFlushInterval = 60 * time.Second // 默认每 60s 触发一次批量处理
)

type BatchFunc[T any] func(t []T) error

type BatchProcessor[T any] struct {
	queue         chan T             // 用于接收数据的队列
	batch         []T                // 当前批次的数据
	flushSignal   chan struct{}      // 用于触发批量处理的信号
	mu            sync.Mutex         // 用于保护批次的互斥锁
	wg            sync.WaitGroup     // 用于等待所有 goroutine 完成
	ctx           context.Context    // 上下文
	cancel        context.CancelFunc // 取消函数
	maxBatchSize  int                // 批次大小
	flushInterval time.Duration      // 刷新间隔
	fn            BatchFunc[T]
}

func NewBatchProcessor[T any](ctx context.Context, maxBatchSize int, flushInterval time.Duration, batchFn BatchFunc[T]) *BatchProcessor[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &BatchProcessor[T]{
		queue:         make(chan T, maxBatchSize),
		batch:         make([]T, 0, maxBatchSize),
		flushSignal:   make(chan struct{}),
		ctx:           ctx,
		cancel:        cancel,
		maxBatchSize:  maxBatchSize,
		flushInterval: flushInterval,
		fn:            batchFn,
	}
}

func (b *BatchProcessor[T]) Start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.run()
	}()
}

func (b *BatchProcessor[T]) Stop() {
	b.cancel() // 取消上下文
	close(b.queue)
	b.wg.Wait()
}

func (b *BatchProcessor[T]) run() {
	ticker := time.NewTicker(b.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case t := <-b.queue:
			b.mu.Lock()
			b.batch = append(b.batch, t)
			if len(b.batch) >= b.maxBatchSize {
				b.flush()
			}
			b.mu.Unlock()
		case <-ticker.C:
			b.flushWithLock()
		case <-b.flushSignal:
			b.flushWithLock()
		case <-b.ctx.Done():
			b.flushWithLock()
			return
		}
	}
}

func (b *BatchProcessor[T]) flush() {
	if err := b.fn(b.batch); err != nil {
		fmt.Printf("Batch processor flush error:%v \n", err)
		// todo 重试
		return
	}
	b.batch = b.batch[:0] // 清空批次
}

func (b *BatchProcessor[T]) flushWithLock() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.batch) > 0 {
		b.flush()
	}
}

func (b *BatchProcessor[T]) Submit(t T) error {
	select {
	case b.queue <- t:
		return nil
	case <-b.ctx.Done():
		return fmt.Errorf("batch processor is stopped")
	}
}
