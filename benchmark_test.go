package benchmark

import (
	"context"
	"testing"
	"time"

	cb "github.com/golang-queue/benchmark/CircularBuffer"
	cl "github.com/golang-queue/benchmark/ContainerList"
	rb "github.com/golang-queue/benchmark/RingBuffer"
	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/queue/job"
)

var count = 1

type testqueue interface {
	Queue(task core.QueuedMessage) error
	Request() (core.QueuedMessage, error)
}

func testDeQueue(b *testing.B, pool testqueue) {
	message := job.NewTask(func(context.Context) error {
		return nil
	},
		job.AllowOption{
			RetryCount: job.Int64(100),
			RetryDelay: job.Time(30 * time.Millisecond),
			Timeout:    job.Time(3 * time.Millisecond),
		},
	)
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			_ = pool.Queue(message)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			_, _ = pool.Request()
		}
	}
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
	b.ResetTimer()
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

func BenchmarkQueueAndRequestContainerList(b *testing.B) {
	pool := cl.NewContainerList(
		b.N * count,
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkQueueAndRequestCircularBuffer(b *testing.B) {
	pool := cb.NewCircularBuffer(
		b.N * count,
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkQueueAndRequestRingBuffer(b *testing.B) {
	pool := rb.NewConsumer(
		b.N * count,
	)

	testQueueAndRequest(b, pool)
}

func BenchmarkQueueCircularBuffer(b *testing.B) {
	pool := cb.NewCircularBuffer(
		b.N * count,
	)

	testQueue(b, pool)
}

func BenchmarkQueueRingBuffer(b *testing.B) {
	pool := rb.NewConsumer(
		b.N * count,
	)

	testQueue(b, pool)
}

func BenchmarkQueueContainerList(b *testing.B) {
	pool := cl.NewContainerList(
		b.N * count,
	)

	testQueue(b, pool)
}

func BenchmarkDeQueueWithRingBuffer(b *testing.B) {
	pool := rb.NewConsumer(
		b.N * count,
	)

	testDeQueue(b, pool)
}

func BenchmarkDeQueueWithContainerList(b *testing.B) {
	pool := cl.NewContainerList(
		b.N * count,
	)

	testDeQueue(b, pool)
}

func BenchmarkDeQueueWithCircularBuffer(b *testing.B) {
	pool := cb.NewCircularBuffer(
		b.N * count,
	)

	testDeQueue(b, pool)
}
