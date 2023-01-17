package benchmark

import (
	"context"
	"testing"
	"time"

	cb "github.com/golang-queue/benchmark/CircularBuffer"
	cl "github.com/golang-queue/benchmark/ContainerList"
	rb "github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/queue/job"

	"github.com/enriquebris/goconcurrentqueue"
)

var count = 1

type testqueue interface {
	Queue(task core.QueuedMessage) error
	Request() (core.QueuedMessage, error)
}

func testQueue(b *testing.B, pool testqueue) {
	message := job.NewTask(func(context.Context) error {
		return nil
	},
		job.AllowOption{
			RetryCount: job.Int64(100),
			RetryDelay: job.Time(30 * time.Millisecond),
			Timeout:    job.Time(3 * time.Millisecond),
		},
	)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			_ = pool.Queue(message)
		}
	}
}

func testQueueAndRequest(b *testing.B, pool testqueue) {
	message := job.NewTask(func(context.Context) error {
		return nil
	},
		job.AllowOption{
			RetryCount: job.Int64(100),
			RetryDelay: job.Time(30 * time.Millisecond),
			Timeout:    job.Time(3 * time.Millisecond),
		},
	)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			_ = pool.Queue(message)
			_, _ = pool.Request()
		}
	}
}

func BenchmarkContainerListQueueAndRequest(b *testing.B) {
	pool := cl.NewContainerList(
		b.N * count,
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkCircularBufferQueueAndRequest(b *testing.B) {
	pool := cb.NewCircularBuffer(
		b.N * count,
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkRingBufferQueueAndRequest(b *testing.B) {
	pool := rb.NewConsumer(
		rb.WithQueueSize(b.N * count),
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkCircularBufferQueue(b *testing.B) {
	pool := cb.NewCircularBuffer(
		b.N * count,
	)

	testQueue(b, pool)
}

func BenchmarkRingBufferQueue(b *testing.B) {
	pool := rb.NewConsumer(
		rb.WithQueueSize(b.N * count),
	)

	testQueue(b, pool)
}

func BenchmarkContainerListQueue(b *testing.B) {
	pool := cl.NewContainerList(
		b.N * count,
	)

	testQueue(b, pool)
}

func BenchmarkFIFOEnqueueSingleGR(b *testing.B) {
	fifo := goconcurrentqueue.NewFIFO()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fifo.Enqueue(i)
		fifo.Dequeue()
	}
}

func BenchmarkFixedFIFOEnqueueSingleGR(b *testing.B) {
	fifo := goconcurrentqueue.NewFixedFIFO(5)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fifo.Enqueue(i)
		fifo.Dequeue()
	}
}
