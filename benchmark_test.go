package benchmark

import (
	"context"
	"testing"
	"time"

	ch "github.com/golang-queue/benchmark/Channel"
	cb "github.com/golang-queue/benchmark/CircularBuffer"
	dl "github.com/golang-queue/benchmark/DoublyLinked"
	rb "github.com/golang-queue/benchmark/RingBuffer"

	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/queue/job"
)

var (
	count  = 1
	result core.QueuedMessage
)

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

	var m core.QueuedMessage
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			m, _ = pool.Request()
		}
	}
	result = m
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

	var m core.QueuedMessage
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			_ = pool.Queue(message)
			m, _ = pool.Request()
		}
	}
	result = m
}

func BenchmarkQueueAndRequestChannel(b *testing.B) {
	testQueueAndRequest(b, ch.NewConsumer(b.N*count))
}

func BenchmarkQueueAndRequestDoublyLinked(b *testing.B) {
	testQueueAndRequest(b, dl.NewDoublyLinked(b.N*count))
}

func BenchmarkQueueAndRequestCircularBuffer(b *testing.B) {
	testQueueAndRequest(b, cb.NewCircularBuffer(b.N*count))
}

func BenchmarkQueueAndRequestRingBuffer(b *testing.B) {
	testQueueAndRequest(b, rb.NewConsumer(b.N*count))
}

func BenchmarkQueueChannel(b *testing.B) {
	testQueue(b, ch.NewConsumer(b.N*count))
}

func BenchmarkQueueCircularBuffer(b *testing.B) {
	testQueue(b, cb.NewCircularBuffer(b.N*count))
}

func BenchmarkQueueRingBuffer(b *testing.B) {
	testQueue(b, rb.NewConsumer(b.N*count))
}

func BenchmarkQueueDoublyLinked(b *testing.B) {
	testQueue(b, dl.NewDoublyLinked(b.N*count))
}

func BenchmarkDeQueueWithChannel(b *testing.B) {
	testDeQueue(b, ch.NewConsumer(b.N*count))
}

func BenchmarkDeQueueWithRingBuffer(b *testing.B) {
	testDeQueue(b, rb.NewConsumer(b.N*count))
}

func BenchmarkDeQueueWithDoublyLinked(b *testing.B) {
	testDeQueue(b, dl.NewDoublyLinked(b.N*count))
}

func BenchmarkDeQueueWithCircularBuffer(b *testing.B) {
	testDeQueue(b, cb.NewCircularBuffer(b.N*count))
}
