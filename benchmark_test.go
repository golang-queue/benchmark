package benchmark

import (
	"context"
	"testing"
	"time"

	q1 "github.com/golang-queue/benchmark/CircularBuffer"
	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/queue/job"

	"github.com/enriquebris/goconcurrentqueue"
)

var count = 1

type testqueue interface {
	Queue(task core.QueuedMessage) error
	Request() (core.QueuedMessage, error)
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
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			_ = pool.Queue(message)
			_, _ = pool.Request()
		}
	}
}

func BenchmarkCircularBuffer(b *testing.B) {
	pool := q1.NewCircularBuffer(
		b.N * count,
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkRingBuffer(b *testing.B) {
	pool := queue.NewConsumer(
		queue.WithQueueSize(b.N * count),
	)

	testQueueAndRequest(b, pool)
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
