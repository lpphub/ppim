package ext

import "sync"

type SegmentRWLock struct {
	segments []sync.RWMutex
	size     uint64
}

func NewSegmentLock(size uint64) *SegmentRWLock {
	sl := &SegmentRWLock{
		segments: make([]sync.RWMutex, size),
		size:     size,
	}
	return sl
}

func (sl *SegmentRWLock) Lock(index uint64) {
	sl.getSegmentLock(index).Lock()
}

func (sl *SegmentRWLock) Unlock(index uint64) {
	sl.getSegmentLock(index).Unlock()
}

func (sl *SegmentRWLock) RLock(index uint64) {
	sl.getSegmentLock(index).RLock()
}

func (sl *SegmentRWLock) RUnlock(index uint64) {
	sl.getSegmentLock(index).RUnlock()
}

func (sl *SegmentRWLock) getSegmentLock(index uint64) *sync.RWMutex {
	return &sl.segments[index%sl.size]
}
