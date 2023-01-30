package benchmark

import (
	"testing"
	"time"

	"github.com/golang-queue/queue/job"

	"github.com/goccy/go-json"
)

type mockMessage struct {
	message string
}

func (m mockMessage) Bytes() []byte {
	return []byte(m.message)
}

func BenchmarkEncode(b *testing.B) {
	m := job.NewMessage(&mockMessage{
		message: "foo",
	}, job.AllowOption{
		RetryCount: job.Int64(100),
		RetryDelay: job.Time(30 * time.Millisecond),
		Timeout:    job.Time(3 * time.Millisecond),
	})

	b.Run("JSON", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(m)
		}
	})

	b.Run("UnsafeCast", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = job.Encode(m)
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	m := job.NewMessage(&mockMessage{
		message: "foo",
	}, job.AllowOption{
		RetryCount: job.Int64(100),
		RetryDelay: job.Time(30 * time.Millisecond),
		Timeout:    job.Time(3 * time.Millisecond),
	})

	b.Run("JSON", func(b *testing.B) {
		data, _ := json.Marshal(m)
		var out *job.Message
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = json.Unmarshal(data, out)
		}
	})

	b.Run("UnsafeCast", func(b *testing.B) {
		data := job.Encode(m)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = job.Decode(data)
		}
	})
}
