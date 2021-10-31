package queue

import (
	"sync"

	"github.com/pkg/errors"
)

var ErrEmptyQueue = errors.New("queue is empty")

// queue is a concurrency-safe fixed size queue of items.
// Implements the Queue interface.
// note: items could have specified type after go 1.18 since the queue could be generic (hold any underlying data structure)
type fixedQueue struct {
	items   []interface{}
	maxSize int
	lock    sync.RWMutex
}

// NewFixedQueue creates a new queue with the given max size
func NewFixedQueue(size int) *fixedQueue {
	return &fixedQueue{
		items:   make([]interface{}, 0, size),
		maxSize: size,
	}
}

// Enqueue adds a item to the queue.
// If queue has reached max capacity, it will remove the oldest item.
func (q *fixedQueue) Enqueue(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == q.maxSize {
		q.items = append(q.items[1:], item)
	} else {
		q.items = append(q.items, item)
	}
}

// Dequeue removes the oldest item from the queue.
func (q *fixedQueue) Dequeue() (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return nil, ErrEmptyQueue
	}

	first, rest := q.items[0], q.items[1:]
	q.items = rest
	return first, nil
}

// IsEmpty returns true if the queue is empty.
func (q *fixedQueue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return len(q.items) == 0
}

// IsFull returns true if the queue is full.
func (q *fixedQueue) IsFull() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return len(q.items) == q.maxSize
}

// Head returns the oldest item in the queue.
func (q *fixedQueue) Head() (interface{}, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	if len(q.items) == 0 {
		return nil, ErrEmptyQueue
	}

	return q.items[0], nil
}
