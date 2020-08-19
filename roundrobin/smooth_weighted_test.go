package roundrobin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSW_Next(t *testing.T) {
	w := NewSmoothRoundrobin()

	s, _ := w.Next()
	assert.Nil(t, s)

	w.Add("server1", 5)
	s, _ = w.Next()
	assert.Equal(t, "server1", s.(string))

	w.Reset()
	s, _ = w.Next()
	assert.Nil(t, s)

	w.Add("server1", 5)
	s, _ = w.Next()
	assert.Equal(t, "server1", s.(string))

	w.Add("server2", 2)
	w.Add("server3", 3)

	results := make(map[string]int)

	for i := 0; i < 1000; i++ {
		s, _ := w.Next()
		results[s.(string)]++
	}

	if results["server1"] != 500 || results["server2"] != 200 || results["server3"] != 300 {
		t.Error("the algorithm is wrong")
	}

	w.(*smoothRoundrobin).items[0].EffectiveWeight = w.(*smoothRoundrobin).items[0].CurrentWeight - 1
	s, _ = w.Next()
	assert.Equal(t, "server3", s.(string))
}
