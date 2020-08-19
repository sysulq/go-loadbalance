package aperture

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
)

func TestAperture(t *testing.T) {
	t.Run("0 item", func(t *testing.T) {
		ll := NewPeakEwmaAperture()
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Nil(t, item)
	})

	t.Run("1 client 1 server", func(t *testing.T) {
		ll := NewSmoothRoundrobin()
		ll.SetLocalPeers(nil)
		ll.SetLocalPeers([]string{"1"})
		ll.SetRemotePeers([]interface{}{"8"})
		ll.SetLocalPeerID("1")

		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "8", item)
	})

	t.Run("1 client 1 server", func(t *testing.T) {
		ll := NewLeastLoadedApeture()
		ll.SetLocalPeers(nil)
		ll.SetLocalPeers([]string{"1"})
		ll.SetRemotePeers([]interface{}{"8"})
		ll.SetLocalPeerID("1")

		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "8", item)
	})

	t.Run("3 client 3 server", func(t *testing.T) {
		ll := NewLeastLoadedApeture()
		ll.SetLocalPeers([]string{"1", "2", "3"})
		ll.SetRemotePeers([]interface{}{"8", "9", "10"})
		ll.SetLocalPeerID("1")
		ll.SetLogicalAperture(1)

		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "8", item)

		ll.SetLocalPeerID("2")

		item, done = ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "9", item)

		ll.SetLocalPeerID("3")

		item, done = ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "10", item)
	})

	t.Run("count", func(t *testing.T) {
		ll := NewLeastLoadedApeture()
		ll.SetLocalPeers([]string{"1", "2", "3"})
		ll.SetRemotePeers([]interface{}{"8", "9", "10", "11", "12"})
		ll.SetLocalPeerID("1")
		ll.SetLogicalAperture(2)

		countMap := make(map[interface{}]int)

		totalCount := 5000
		wg := sync.WaitGroup{}
		wg.Add(totalCount)

		mu := sync.Mutex{}
		for i := 0; i < totalCount; i++ {
			go func() {
				defer wg.Done()
				item, done := ll.Next()
				time.Sleep(1 * time.Second)
				done(balancer.DoneInfo{})

				mu.Lock()
				countMap[item]++
				mu.Unlock()
			}()
		}

		wg.Wait()

		total := 0
		for _, count := range countMap {
			total += count
		}
		assert.Less(t, totalCount*3/10-10, countMap["8"])
		assert.Less(t, totalCount*3/10-10, countMap["9"])
		assert.Less(t, totalCount*3/10-10, countMap["10"])
		assert.Less(t, totalCount*1/10-10, countMap["11"])

		assert.Equal(t, totalCount, total)
	})
}

func TestDynamic(t *testing.T) {
	t.Run("1client-3client", func(t *testing.T) {
		ll := NewLeastLoadedApeture()
		ll.SetLocalPeers([]string{"1"})
		ll.SetRemotePeers([]interface{}{"8", "9", "10"})
		ll.SetLocalPeerID("1")
		ll.SetLogicalAperture(2)

		assert.Equal(t, []int{0, 1, 2}, ll.(*aperture).List())

		ll.SetLocalPeers([]string{"1", "2", "3"})
		assert.Equal(t, []int{0, 1}, ll.(*aperture).List())

	})

	t.Run("3server-4server", func(t *testing.T) {
		ll := NewLeastLoadedApeture()
		ll.SetLocalPeers([]string{"1", "2", "3"})
		ll.SetRemotePeers([]interface{}{"8", "9", "10"})
		ll.SetLocalPeerID("1")
		ll.SetLogicalAperture(2)

		assert.Equal(t, []int{0, 1}, ll.(*aperture).List())

		ll.SetRemotePeers([]interface{}{"1", "2", "3", "4"})
		assert.Equal(t, []int{0, 1, 2}, ll.(*aperture).List())

	})
}
