package p2c_test

import (
	"testing"

	"github.com/hnlq715/go-loadbalance/p2c"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
)

func TestLeastLoaded(t *testing.T) {
	t.Run("0 item", func(t *testing.T) {
		ll := p2c.NewLeastLoaded()
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Nil(t, item)
	})

	t.Run("1 item", func(t *testing.T) {
		ll := p2c.NewLeastLoaded()
		ll.Add(1, 1)
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, 1, item)
	})

	t.Run("3 items", func(t *testing.T) {
		ll := p2c.NewLeastLoaded()
		ll.Add(1, 1)
		ll.Add(2, 1)
		ll.Add(3, 1)

		countMap := make(map[interface{}]int)

		totalCount := 10000
		for i := 0; i < totalCount; i++ {
			item, done := ll.Next()
			done(balancer.DoneInfo{})

			countMap[item]++
		}

		total := 0
		for _, count := range countMap {
			total += count
			assert.Less(t, totalCount/3-200, count)
		}

		assert.Equal(t, totalCount, total)
	})
}

func TestLeastLoadedAbnormal(t *testing.T) {
	t.Run("fixed inflight", func(t *testing.T) {
		ll := p2c.NewLeastLoaded()
		ll.Add(1, 1)
		ll.Add(2, 1)
		ll.Add(3, 1)

		item, _ := ll.Next()

		for i := 0; i < 1000; i++ {
			next, done := ll.Next()
			done(balancer.DoneInfo{})

			assert.NotEqual(t, item, next)
		}
	})
}
