package ringbuffer

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

var _ core.Worker = (*Consumer)(nil)

// Consumer for simple queue using buffer channel
type Consumer struct {
	sync.Mutex
	taskQueue []core.TaskMessage
	runFunc   func(context.Context, core.TaskMessage) error
	capacity  int
	count     int
	head      int
	tail      int
	exit      chan struct{}
	stopOnce  sync.Once
	stopFlag  atomic.Int32
}

// Run to execute new task
func (s *Consumer) Run(ctx context.Context, task core.TaskMessage) error {
	return s.runFunc(ctx, task)
}

// Shutdown the worker
func (s *Consumer) Shutdown() error {
	if !s.stopFlag.CompareAndSwap(0, 1) {
		return queue.ErrQueueShutdown
	}

	s.stopOnce.Do(func() {
		if s.count > 0 {
			<-s.exit
		}
	})
	return nil
}

// Queue send task to the buffer channel
func (s *Consumer) Queue(task core.TaskMessage) error {
	if s.stopFlag.Load() == 1 {
		return queue.ErrQueueShutdown
	}
	if s.capacity > 0 && s.count >= s.capacity {
		return queue.ErrMaxCapacity
	}

	s.Lock()
	if s.count == len(s.taskQueue) {
		s.resize(s.count * 2)
	}
	s.taskQueue[s.tail] = task
	s.tail = (s.tail + 1) % len(s.taskQueue)
	s.count++
	s.Unlock()

	return nil
}

// Request a new task from channel
func (s *Consumer) Request() (core.TaskMessage, error) {
	if s.stopFlag.Load() == 1 && s.count == 0 {
		select {
		case s.exit <- struct{}{}:
		default:
		}
		return nil, queue.ErrQueueHasBeenClosed
	}

	if s.count == 0 {
		return nil, queue.ErrNoTaskInQueue
	}
	s.Lock()
	data := s.taskQueue[s.head]
	s.taskQueue[s.head] = nil
	s.head = (s.head + 1) % len(s.taskQueue)
	s.count--

	if n := len(s.taskQueue) / 2; n > 2 && s.count <= n {
		s.resize(n)
	}
	s.Unlock()

	return data, nil
}

func (s *Consumer) resize(n int) {
	nodes := make([]core.TaskMessage, n)
	if s.head < s.tail {
		copy(nodes, s.taskQueue[s.head:s.tail])
	} else {
		copy(nodes, s.taskQueue[s.head:])
		copy(nodes[len(s.taskQueue)-s.head:], s.taskQueue[:s.tail])
	}

	s.tail = s.count % n
	s.head = 0
	s.taskQueue = nodes
}

// NewConsumer for create new Consumer instance
func NewConsumer(size int) *Consumer {
	w := &Consumer{
		taskQueue: make([]core.TaskMessage, 2),
		capacity:  size,
		exit:      make(chan struct{}),
	}

	return w
}
