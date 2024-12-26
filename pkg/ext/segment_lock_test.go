package ext

import (
	"testing"
	"time"
)

func TestSegmentLock_Lock(t *testing.T) {
	sl := NewSegmentLock(2)

	for i := 0; i < 10; i++ {
		go func() {
			sl.Lock(uint32(i))
			defer sl.Unlock(uint32(i))

			t.Logf("lock %d", i)
		}()
	}

	time.Sleep(10 * time.Second)
}
