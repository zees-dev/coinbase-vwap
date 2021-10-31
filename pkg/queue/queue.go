package queue

// Queue is an interface for a generic queue implementation
type Queue interface {
	Enqueue(interface{})
	Dequeue() (interface{}, error)
	IsEmpty() bool
	IsFull() bool
	Head() (interface{}, error)
}
