package benchmark

import (
	"testing"
	"time"

	"github.com/golang-queue/queue/job"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
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
		out := &job.Message{}
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

func TestEncodeAndDecode(t *testing.T) {
	m := job.NewMessage(&mockMessage{
		message: "foo",
	}, job.AllowOption{
		RetryCount: job.Int64(100),
		RetryDelay: job.Time(30 * time.Millisecond),
		Timeout:    job.Time(3 * time.Millisecond),
	})

	t.Run("JSON", func(t *testing.T) {
		data, _ := json.Marshal(m)
		out := &job.Message{}
		_ = json.Unmarshal(data, out)

		assert.Equal(t, int64(100), out.RetryCount)
		assert.Equal(t, 30*time.Millisecond, out.RetryDelay)
		assert.Equal(t, 3*time.Millisecond, out.Timeout)
	})

	t.Run("UnsafeCast", func(t *testing.T) {
		data := job.Encode(m)
		out := job.Decode(data)

		assert.Equal(t, int64(100), out.RetryCount)
		assert.Equal(t, 30*time.Millisecond, out.RetryDelay)
		assert.Equal(t, 3*time.Millisecond, out.Timeout)
	})
}
