package ext

import (
	"fmt"
	"testing"
	"time"
)

func TestNewBatchProcessor(t *testing.T) {
	type Message struct {
		ID      int
		Content string
	}

	bp := NewBatchProcessor[Message](2, 10, 3*time.Second, func(ms []Message) error {
		// 处理消息的逻辑
		//t.Logf("Processing %d messages\n", len(m))
		for _, m := range ms {
			t.Logf("Processing message: %+v\n", m)
		}
		return nil
	})
	bp.Start()

	// 模拟消息生产
	go func() {
		for i := 1; i <= 25; i++ {
			msg := Message{
				ID:      i,
				Content: fmt.Sprintf("Message %d", i),
			}
			if err := bp.Submit(msg); err != nil {
				t.Logf("Failed to add message: %v\n", err)
				break
			}
			time.Sleep(2 * time.Millisecond) // 模拟消息间隔
		}
	}()

	// 模拟程序运行一段时间后退出
	time.Sleep(10 * time.Second)
	bp.Stop()
}

func TestNewGroupedBatchProcessor(t *testing.T) {
	fn := func(str []string) error {
		fmt.Printf("print: %v \n", str)
		return nil
	}

	batchList := []*BatchProcessor[string]{
		NewBatchProcessor[string](2, 2, 3*time.Second, fn),
		NewBatchProcessor[string](1, 3, 3*time.Second, fn),
	}

	group := NewGroupedBatchProcessor(batchList)
	group.Start()
	defer group.Stop()

	go func() {
		for i := 1; i <= 100; i++ {
			msg := fmt.Sprintf("test-%d", i)
			_ = group.Submit(msg, msg)

			time.Sleep(10 * time.Millisecond)
		}
	}()

	time.Sleep(10 * time.Second)
}
