package ext

import (
	"errors"
)

type node[T any] struct {
	val  T
	next *node[T]
}
type LinkedQueue[T any] struct {
	head *node[T]
	tail *node[T]
	size int
}

// Enqueue 入队
func (q *LinkedQueue[T]) Enqueue(value T) {
	newNode := &node[T]{val: value}
	if q.head == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}
	q.size++
}

// Dequeue 出队
func (q *LinkedQueue[T]) Dequeue() (T, error) {
	if q.head == nil {
		var zero T
		return zero, errors.New("queue is empty")
	}
	value := q.head.val
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return value, nil
}

func (q *LinkedQueue[T]) Peek() (T, error) {
	if q.head == nil {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.head.val, nil
}

func (q *LinkedQueue[T]) Take(num int) ([]T, error) {
	if num <= 0 {
		return nil, errors.New("num must be greater than 0")
	}
	if q.size == 0 {
		return nil, errors.New("queue is empty")
	}

	result := make([]T, 0, num)
	for i := 0; i < num && q.head != nil; i++ {
		value := q.head.val
		q.head = q.head.next
		result = append(result, value)
		q.size--
	}

	if q.head == nil {
		q.tail = nil
	}
	return result, nil
}

func (q *LinkedQueue[T]) Size() int {
	return q.size
}
