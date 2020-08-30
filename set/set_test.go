package set_test

import (
	"testing"

	"github.com/hnlq715/go-loadbalance"
	"github.com/hnlq715/go-loadbalance/set"
	"gotest.tools/assert"
)

func TestSet(t *testing.T) {
	t.Run("same set", func(t *testing.T) {
		s := set.New(loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "01",
		})

		s.Add(1, 1, loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "01",
		})

		item, _ := s.Next()
		assert.Equal(t, 1, item)

		s.Reset()

		item, _ = s.Next()
		assert.Equal(t, nil, item)
	})

	t.Run("different region", func(t *testing.T) {
		s := set.New(loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "01",
		})

		s.Add(1, 1, loadbalance.SetInfo{
			Name:     "app",
			Region:   "sh",
			UnitName: "02",
		})

		item, _ := s.Next()
		assert.Equal(t, nil, item)
	})

	t.Run("different name", func(t *testing.T) {
		s := set.New(loadbalance.SetInfo{
			Name:     "app01",
			Region:   "bj",
			UnitName: "01",
		})

		s.Add(1, 1, loadbalance.SetInfo{
			Name:     "app02",
			Region:   "bj",
			UnitName: "01",
		})

		item, _ := s.Next()
		assert.Equal(t, nil, item)
	})

	t.Run("different set", func(t *testing.T) {
		s := set.New(loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "01",
		})

		s.Add(1, 1, loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "02",
		})

		item, _ := s.Next()
		assert.Equal(t, nil, item)
	})

	t.Run("* set", func(t *testing.T) {
		s := set.New(loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "*",
		})

		s.Add(1, 1, loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "*",
		})

		s.Add(2, 1, loadbalance.SetInfo{
			Name:     "app",
			Region:   "bj",
			UnitName: "01",
		})

		item, _ := s.Next()
		assert.Equal(t, 1, item)

		item, _ = s.Next()
		assert.Equal(t, 2, item)
	})
}
