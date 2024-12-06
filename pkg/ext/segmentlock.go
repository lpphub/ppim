package ext

import "sync"

type SegmentLock struct {
	segments []sync.RWMutex
	size     uint32
}

func NewSegmentLock(size uint32) *SegmentLock {
	sl := &SegmentLock{
		segments: make([]sync.RWMutex, size),
		size:     size,
	}
	return sl
}

func (sl *SegmentLock) Lock(index uint32) {
	sl.getSegmentLock(index).Lock()
}

func (sl *SegmentLock) Unlock(index uint32) {
	sl.getSegmentLock(index).Unlock()
}

func (sl *SegmentLock) RLock(index uint32) {
	sl.getSegmentLock(index).RLock()
}

func (sl *SegmentLock) RUnlock(index uint32) {
	sl.getSegmentLock(index).RUnlock()
}

func (sl *SegmentLock) getSegmentLock(index uint32) *sync.RWMutex {
	return &sl.segments[index%sl.size]
}
