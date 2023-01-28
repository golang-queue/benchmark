package circularbuffer

import (
	"sync"
	"sync/atomic"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

// CircularBuffer for simple queue using buffer channel
type CircularBuffer struct {
	sync.Mutex
	taskQueue []core.QueuedMessage
	capacity  int
	head      int
	tail      int
	exit      chan struct{}
	stopOnce  sync.Once
	stopFlag  int32
}

// Shutdown the worker
func (s *CircularBuffer) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&s.stopFlag, 0, 1) {
		return queue.ErrQueueShutdown
	}

	s.stopOnce.Do(func() {
		if !s.IsEmpty() {
			<-s.exit
		}
	})
	return nil
}

// Queue send task to the buffer channel
func (s *CircularBuffer) Queue(task core.QueuedMessage) error {
	if atomic.LoadInt32(&s.stopFlag) == 1 {
		return queue.ErrQueueShutdown
	}
	if s.IsFull() {
		return queue.ErrMaxCapacity
	}

	s.Lock()
	s.taskQueue[s.tail] = task
	s.tail = (s.tail + 1) % s.capacity
	s.Unlock()

	return nil
}

// Request a new task from channel
func (s *CircularBuffer) Request() (core.QueuedMessage, error) {
	if atomic.LoadInt32(&s.stopFlag) == 1 && s.IsEmpty() {
		select {
		case s.exit <- struct{}{}:
		default:
		}
		return nil, queue.ErrQueueHasBeenClosed
	}

	if s.IsEmpty() {
		return nil, queue.ErrNoTaskInQueue
	}

	s.Lock()
	data := s.taskQueue[s.head]
	s.head = (s.head + 1) % s.capacity
	s.Unlock()

	return data, nil
}

// IsEmpty returns true if queue is empty
func (s *CircularBuffer) IsEmpty() bool {
	return s.head == s.tail
}

// IsFull returns true if queue is full
func (s *CircularBuffer) IsFull() bool {
	return s.head == (s.tail+1)%s.capacity
}

// NewCircularBuffer for create new CircularBuffer instance
func NewCircularBuffer(size int) *CircularBuffer {
	w := &CircularBuffer{
		taskQueue: make([]core.QueuedMessage, size),
		capacity:  size,
		exit:      make(chan struct{}),
	}

	return w
}
