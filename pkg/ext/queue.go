package ext

import (
	"errors"
	"sync"
)

type node[T any] struct {
	Val  T
	Next *node[T]
}
type LinkedListQueue[T any] struct {
	head *node[T]
	tail *node[T]
	size int
	mtx  sync.Mutex
}

// Enqueue 入队
func (q *LinkedListQueue[T]) Enqueue(value T) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	newNode := &node[T]{Val: value}
	if q.tail != nil {
		q.tail.Next = newNode
	}
	q.tail = newNode
	if q.head == nil {
		q.head = q.tail
	}
	q.size++
}

// Dequeue 出队
func (q *LinkedListQueue[T]) Dequeue() (T, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.head == nil {
		var zero T
		return zero, errors.New("queue is empty")
	}
	value := q.head.Val
	q.head = q.head.Next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return value, nil
}

func (q *LinkedListQueue[T]) Peek() (T, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.head == nil {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.head.Val, nil
}

func (q *LinkedListQueue[T]) Take(num int) ([]T, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if num <= 0 {
		return nil, errors.New("num must be greater than 0")
	}
	if q.size == 0 {
		return nil, errors.New("queue is empty")
	}

	var result []T
	for i := 0; i < num && q.head != nil; i++ {
		value := q.head.Val
		q.head = q.head.Next
		result = append(result, value)
		q.size--
	}

	if q.head == nil {
		q.tail = nil
	}
	return result, nil
}

func (q *LinkedListQueue[T]) Size() int {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.size
}
