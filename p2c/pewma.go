package p2c

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hnlq715/go-loadbalance"
	"google.golang.org/grpc/balancer"
)

type peakEwma struct {
	stamp int64
	value int64
	tau   time.Duration
}

const (
	defaultTau = 10000 * time.Millisecond
)

func newPEWMA() *peakEwma {
	return &peakEwma{
		tau: defaultTau,
	}
}

// Observe 计算peak指数加权移动平均值
func (p *peakEwma) Observe(rtt int64) {
	now := time.Now().UnixNano()

	stamp := atomic.SwapInt64(&p.stamp, now)
	td := now - stamp

	if td < 0 {
		td = 0
	}

	w := math.Exp(float64(-td) / float64(p.tau))
	latency := atomic.LoadInt64(&p.value)

	if rtt > latency {
		atomic.StoreInt64(&p.value, rtt)
	} else {
		atomic.StoreInt64(&p.value, int64(float64(latency)*w+float64(rtt)*(1.0-w)))
	}
}

func (p *peakEwma) Value() int64 {
	return atomic.LoadInt64(&p.value)
}

type peakEwmaNode struct {
	item    interface{}
	latency *peakEwma
	weight  float64
}

type pewma struct {
	items []*peakEwmaNode
	mu    sync.Mutex
	rand  *rand.Rand
}

func NewPeakEwma() loadbalance.P2C {
	return &pewma{
		items: make([]*peakEwmaNode, 0),
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *pewma) Add(item interface{}, weight float64) {
	p.items = append(p.items, &peakEwmaNode{item: item, latency: newPEWMA(), weight: weight})
}

func (p *pewma) Next() (interface{}, func(balancer.DoneInfo)) {
	var sc, backsc *peakEwmaNode
	begin := time.Now().UnixNano()

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
		if float64(sc.latency.Value())*backsc.weight > float64(backsc.latency.Value())*sc.weight {
			sc, backsc = backsc, sc
		}
	}

	return sc.item, func(balancer.DoneInfo) {
		end := time.Now().UnixNano()
		sc.latency.Observe(end - begin)
	}
}
