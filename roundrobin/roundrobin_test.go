package roundrobin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSW_Next(t *testing.T) {
	w := NewSmoothRoundrobin()

	assert.Nil(t, w.Next())

	w.Add("server1", 5)
	assert.Equal(t, "server1", w.Next().(string))

	w.Add("server2", 2)
	w.Add("server3", 3)

	results := make(map[string]int)

	for i := 0; i < 1000; i++ {
		s := w.Next().(string)
		results[s]++
	}

	if results["server1"] != 500 || results["server2"] != 200 || results["server3"] != 300 {
		t.Error("the algorithm is wrong")
	}
}
