package roundrobin

import (
	"math"

	"github.com/hnlq715/go-loadbalance"
	"github.com/hnlq715/go-loadbalance/internal"
	"google.golang.org/grpc/balancer"
)

// smoothRoundrobinNode is a wrapped weighted item.
type smoothRoundrobinNode struct {
	Item            interface{}
	Weight          int64
	CurrentWeight   int64
	EffectiveWeight int64
}

type smoothRoundrobin struct {
	items []*smoothRoundrobinNode
	n     int64
}

// NewSmoothRoundrobin (Smooth Weighted) contains weighted items and provides methods to select a weighted item.
// It is used for the smooth weighted round-robin balancing algorithm.
// This algorithm is implemented in Nginx:
// https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35.
//
// Algorithm is as follows: on each peer selection we increase current_weight
// of each eligible peer by its weight, select peer with greatest current_weight
// and reduce its current_weight by total number of weight points distributed
// among peers.
// In case of { 5, 1, 1 } weights this gives the following sequence of
// current_weight's: (a, a, b, a, c, a, a)
func NewSmoothRoundrobin() loadbalance.Picker {
	return &smoothRoundrobin{}
}

// Add a weighted server.
func (w *smoothRoundrobin) Add(item interface{}, weight float64) {
	wt := int64(math.Floor(weight))
	weighted := &smoothRoundrobinNode{Item: item, Weight: wt, EffectiveWeight: wt}
	w.items = append(w.items, weighted)
	w.n++
}

func (w *smoothRoundrobin) Reset() {
	w.items = w.items[:0]
	w.n = 0
}

// Next returns next selected server.
func (w *smoothRoundrobin) Next() (interface{}, func(balancer.DoneInfo)) {
	if w.n == 0 {
		return nil, internal.EmptyDoneFunc
	}

	if w.n == 1 {
		return w.items[0].Item, internal.EmptyDoneFunc
	}

	return nextSmoothWeighted(w.items).Item, internal.EmptyDoneFunc
}

// nextSmoothWeighted selects the best node through the smooth weighted roundrobin .
func nextSmoothWeighted(items []*smoothRoundrobinNode) (best *smoothRoundrobinNode) {
	total := int64(0)

	for i := 0; i < len(items); i++ {
		w := items[i]

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight

		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	best.CurrentWeight -= total

	return best
}
