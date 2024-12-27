package ext

import "sync"

type SegmentRWLock struct {
	segments []sync.RWMutex
	size     int
}

func NewSegmentLock(size int) *SegmentRWLock {
	sl := &SegmentRWLock{
		segments: make([]sync.RWMutex, size),
		size:     size,
	}
	return sl
}

func (sl *SegmentRWLock) Lock(index int) {
	sl.getSegmentLock(index).Lock()
}

func (sl *SegmentRWLock) Unlock(index int) {
	sl.getSegmentLock(index).Unlock()
}

func (sl *SegmentRWLock) RLock(index int) {
	sl.getSegmentLock(index).RLock()
}

func (sl *SegmentRWLock) RUnlock(index int) {
	sl.getSegmentLock(index).RUnlock()
}

func (sl *SegmentRWLock) getSegmentLock(index int) *sync.RWMutex {
	return &sl.segments[index%sl.size]
}
