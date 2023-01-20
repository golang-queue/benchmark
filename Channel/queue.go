package queue

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

var _ core.Worker = (*Consumer)(nil)

var errMaxCapacity = errors.New("max capacity reached")

// Consumer for simple queue using buffer channel
type Consumer struct {
	taskQueue chan core.QueuedMessage
	runFunc   func(context.Context, core.QueuedMessage) error
	stop      chan struct{}
	exit      chan struct{}
	stopOnce  sync.Once
	stopFlag  int32
}

// Run to execute new task
func (s *Consumer) Run(ctx context.Context, task core.QueuedMessage) error {
	return s.runFunc(ctx, task)
}

// Shutdown the worker
func (s *Consumer) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&s.stopFlag, 0, 1) {
		return queue.ErrQueueShutdown
	}

	s.stopOnce.Do(func() {
		close(s.stop)
		close(s.taskQueue)
		if len(s.taskQueue) > 0 {
			<-s.exit
		}
	})
	return nil
}

// Queue send task to the buffer channel
func (s *Consumer) Queue(task core.QueuedMessage) error {
	if atomic.LoadInt32(&s.stopFlag) == 1 {
		return queue.ErrQueueShutdown
	}

	select {
	case s.taskQueue <- task:
		return nil
	default:
		return errMaxCapacity
	}
}

// Request a new task from channel
func (s *Consumer) Request() (core.QueuedMessage, error) {
	select {
	case task, ok := <-s.taskQueue:
		if !ok {
			select {
			case s.exit <- struct{}{}:
			default:
			}
			return nil, queue.ErrQueueHasBeenClosed
		}
		return task, nil
	default:
		return nil, queue.ErrNoTaskInQueue
	}
}

// NewConsumer for create new consumer instance
func NewConsumer(size int) *Consumer {
	w := &Consumer{
		taskQueue: make(chan core.QueuedMessage, size),
		stop:      make(chan struct{}),
		exit:      make(chan struct{}),
	}

	return w
}
