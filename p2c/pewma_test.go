package p2c

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
)

func TestPEWMA(t *testing.T) {
	p := newPEWMA()
	// p.tau = 600 * time.Millisecond
	p.Observe(int64(1 * time.Second))
	assert.Equal(t, p.Value(), int64(1*time.Second))
	p.Observe(int64(1 * time.Second))
	assert.Equal(t, p.Value(), int64(1*time.Second))
	p.Observe(int64(1 * time.Second))
	assert.Equal(t, p.Value(), int64(1*time.Second))
	time.Sleep(1 * time.Second)
	p.Observe(int64(1 * time.Second))
	assert.Equal(t, p.Value(), int64(1*time.Second))
	p.Observe(int64(2 * time.Second))
	assert.Equal(t, p.Value(), int64(2*time.Second))
	for i := 0; i <= 1000; i++ {
		time.Sleep(1 * time.Microsecond)
		p.Observe(int64(1 * time.Second))
	}
	assert.True(t, p.Value() > int64(1800*time.Millisecond) && p.Value() < int64(2000*time.Millisecond), fmt.Sprintf("%d", p.Value()))
}

func BenchmarkPeakEwma(b *testing.B) {
	b.ResetTimer()
	p := newPEWMA()
	for i := 0; i < b.N; i++ {
		p.Observe(int64(time.Duration(rand.Intn(10)) * time.Second))
	}
	// b.Error(p.Value())
}

func TestPeakEwma(t *testing.T) {
	t.Run("0 item", func(t *testing.T) {
		ll := NewPeakEwma()
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Nil(t, item)
	})

	t.Run("1 item", func(t *testing.T) {
		ll := NewPeakEwma()
		ll.Reset()
		ll.Add(1, 1)
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, 1, item)
	})

	t.Run("3 items", func(t *testing.T) {
		ll := NewPeakEwma()
		ll.Add(1, 1)
		ll.Add(2, 1)
		ll.Add(3, 1)

		countMap := make(map[interface{}]int)

		totalCount := 1000
		for i := 0; i < totalCount; i++ {
			item, done := ll.Next()
			time.Sleep(time.Millisecond)
			done(balancer.DoneInfo{})

			countMap[item]++
		}

		total := 0
		for _, count := range countMap {
			total += count
			assert.Less(t, totalCount/3-2000, count)
		}

		assert.Equal(t, totalCount, total)
	})
}
