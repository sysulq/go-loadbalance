package p2c

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hnlq715/go-loadbalance"
	"google.golang.org/grpc/balancer"
)

type leastLoadedNode struct {
	item     interface{}
	inflight int64
	weight   float64
}

type leastLoaded struct {
	items []*leastLoadedNode
	mu    sync.Mutex
	rand  *rand.Rand
}

func NewLeastLoaded() loadbalance.Picker {
	return &leastLoaded{
		items: make([]*leastLoadedNode, 0),
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *leastLoaded) Add(item interface{}, weight float64) {
	p.items = append(p.items, &leastLoadedNode{item: item, weight: weight})
}

func (p *leastLoaded) Reset() {
	p.items = p.items[:0]
}

func (p *leastLoaded) Next() (interface{}, func(balancer.DoneInfo)) {
	var sc, backsc *leastLoadedNode

	switch len(p.items) {
	case 0:
		return nil, func(balancer.DoneInfo) {}
	case 1:
		sc = p.items[0]
	default:
		// rand needs lock
		p.mu.Lock()
		a := p.rand.Intn(len(p.items))
		b := p.rand.Intn(len(p.items) - 1)
		p.mu.Unlock()

		if b >= a {
			b++
		}

		sc, backsc = p.items[a], p.items[b]

		// choose the least loaded item based on inflight and weight
		scInflight := atomic.LoadInt64(&sc.inflight)
		backscInflight := atomic.LoadInt64(&backsc.inflight)

		if float64(scInflight)*backsc.weight > float64(backscInflight)*sc.weight {
			sc, backsc = backsc, sc
		}
	}

	atomic.AddInt64(&sc.inflight, 1)

	return sc.item, func(balancer.DoneInfo) {
		atomic.AddInt64(&sc.inflight, -1)
	}
}
