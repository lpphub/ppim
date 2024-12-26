package ext

import "sync"

type SegmentRWLock struct {
	segments []sync.RWMutex
	size     uint32
}

func NewSegmentLock(size uint32) *SegmentRWLock {
	sl := &SegmentRWLock{
		segments: make([]sync.RWMutex, size),
		size:     size,
	}
	return sl
}

func (sl *SegmentRWLock) Lock(index uint32) {
	sl.getSegmentLock(index).Lock()
}

func (sl *SegmentRWLock) Unlock(index uint32) {
	sl.getSegmentLock(index).Unlock()
}

func (sl *SegmentRWLock) RLock(index uint32) {
	sl.getSegmentLock(index).RLock()
}

func (sl *SegmentRWLock) RUnlock(index uint32) {
	sl.getSegmentLock(index).RUnlock()
}

func (sl *SegmentRWLock) getSegmentLock(index uint32) *sync.RWMutex {
	return &sl.segments[index%sl.size]
}
