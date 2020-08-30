package set

import (
	"github.com/hnlq715/go-loadbalance"
	"github.com/hnlq715/go-loadbalance/roundrobin"
	"google.golang.org/grpc/balancer"
)

type Set struct {
	info   loadbalance.SetInfo
	picker loadbalance.Picker
}

func New(info loadbalance.SetInfo) loadbalance.Set {
	return &Set{
		info:   info,
		picker: roundrobin.NewSmoothRoundrobin(),
	}
}

func (s *Set) Next() (interface{}, func(balancer.DoneInfo)) {
	return s.picker.Next()
}

func (s *Set) Add(item interface{}, weigth float64, info loadbalance.SetInfo) {
	if info.Name != s.info.Name {
		return
	}

	if info.Region != s.info.Region {
		return
	}

	if s.info.UnitName != "*" {
		if info.UnitName != s.info.UnitName {
			return
		}
	}

	s.picker.Add(item, weigth)
}

func (s *Set) Reset() {
	s.picker.Reset()
}
