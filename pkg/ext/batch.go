package ext

import (
	"context"
	"fmt"
	"ppim/pkg/util"
	"sync"
	"time"
)

type BatchFunc[T any] func(t []T) error

type BatchProcessor[T any] struct {
	queue         chan T             // 用于接收数据的队列
	batchSize     int                // 批次大小
	workerCount   int                // 协程数量
	flushInterval time.Duration      // 刷新间隔
	process       BatchFunc[T]       // 批量处理函数
	wg            sync.WaitGroup     // 用于等待所有 goroutine 完成
	ctx           context.Context    // 上下文
	cancel        context.CancelFunc // 取消函数
}

// NewBatchProcessor
// 同一批次顺序处理时，workerCount 应为 1；
// 同一批次并发处理时，要确保 batchFn 是线程安全的。
func NewBatchProcessor[T any](batchSize, workerCount int, interval time.Duration, batchFn BatchFunc[T]) *BatchProcessor[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &BatchProcessor[T]{
		queue:         make(chan T, batchSize),
		batchSize:     batchSize,
		workerCount:   workerCount,
		flushInterval: interval,
		process:       batchFn,
		ctx:           ctx,
		cancel:        cancel,
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
			if len(batch) >= p.batchSize {
				p.flush(batch)
				batch = nil
			}
		case <-ticker.C:
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
	if err := p.process(batch); err != nil {
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

type GroupedBatchProcessor[T any] struct {
	processors []*BatchProcessor[T]
}

func NewGroupedBatchProcessor[T any](processors []*BatchProcessor[T]) *GroupedBatchProcessor[T] {
	return &GroupedBatchProcessor[T]{
		processors: processors,
	}
}

func (p *GroupedBatchProcessor[T]) Submit(key string, t T) error {
	index := util.XXHash64(key) % uint64(len(p.processors))
	return p.processors[index].Submit(t)
}

func (p *GroupedBatchProcessor[T]) Start() {
	for _, processor := range p.processors {
		processor.Start()
	}
}

func (p *GroupedBatchProcessor[T]) Stop() {
	for _, processor := range p.processors {
		processor.Stop()
	}
}