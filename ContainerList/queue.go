package containerlist

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

var _ core.Worker = (*ContainerList)(nil)

var errMaxCapacity = errors.New("max capacity reached")

// ContainerList for simple queue using buffer channel
type ContainerList struct {
	taskQueue *list.List
	runFunc   func(context.Context, core.QueuedMessage) error
	capacity  int
	exit      chan struct{}
	stopOnce  sync.Once
	stopFlag  int32
}

// Run to execute new task
func (s *ContainerList) Run(ctx context.Context, task core.QueuedMessage) error {
	return s.runFunc(ctx, task)
}

// Shutdown the worker
func (s *ContainerList) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&s.stopFlag, 0, 1) {
		return queue.ErrQueueShutdown
	}

	s.stopOnce.Do(func() {
		if s.taskQueue.Len() > 0 {
			<-s.exit
		}
	})
	return nil
}

// Queue send task to the buffer channel
func (s *ContainerList) Queue(task core.QueuedMessage) error {
	if atomic.LoadInt32(&s.stopFlag) == 1 {
		return queue.ErrQueueShutdown
	}

	if s.taskQueue.Len() >= s.capacity {
		return errMaxCapacity
	}

	s.taskQueue.PushBack(task)

	return nil
}

// Request a new task from channel
func (s *ContainerList) Request() (core.QueuedMessage, error) {
	if atomic.LoadInt32(&s.stopFlag) == 1 && s.taskQueue.Len() == 0 {
		select {
		case s.exit <- struct{}{}:
		default:
		}
		return nil, queue.ErrQueueHasBeenClosed
	}

	if s.taskQueue.Len() == 0 {
		return nil, queue.ErrNoTaskInQueue
	}

	peak := s.taskQueue.Back()
	s.taskQueue.Remove(s.taskQueue.Back())

	return peak.Value.(core.QueuedMessage), nil
}

// NewContainerList for create new ContainerList instance
func NewContainerList(size int) *ContainerList {
	w := &ContainerList{
		taskQueue: list.New(),
		capacity:  size,
		exit:      make(chan struct{}),
	}

	return w
}
