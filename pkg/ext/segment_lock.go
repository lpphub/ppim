package ext

import "sync"

type SegmentLock struct {
	segments []sync.RWMutex
	size     int
}

func NewSegmentLock(size int) *SegmentLock {
	sl := &SegmentLock{
		segments: make([]sync.RWMutex, size),
		size:     size,
	}
	return sl
}

func (sl *SegmentLock) Lock(index int) {
	sl.getSegmentLock(index).Lock()
}

func (sl *SegmentLock) Unlock(index int) {
	sl.getSegmentLock(index).Unlock()
}

func (sl *SegmentLock) RLock(index int) {
	sl.getSegmentLock(index).RLock()
}

func (sl *SegmentLock) RUnlock(index int) {
	sl.getSegmentLock(index).RUnlock()
}

func (sl *SegmentLock) getSegmentLock(index int) *sync.RWMutex {
	index = index % sl.size
	if index < 0 {
		index += sl.size
	}
	return &sl.segments[index]
}
