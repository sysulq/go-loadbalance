package aperture

import (
	"math"
)

type Ring struct {
	size      int
	unitWidth float64
}

const (
	floatOne float64 = 1.0
	intOne   int     = 1
)

func NewRing(size int) *Ring {
	return &Ring{
		size:      size,
		unitWidth: floatOne / float64(size),
	}
}

func (r *Ring) Range(offset, width float64) int {
	begin := r.Index(offset)
	end := r.Index(math.Mod(offset+width, 1.0))

	if begin == end && width > r.unitWidth {
		return r.size
	} else if begin == end {
		return intOne
	}

	beginWeight := r.Weight(begin, offset, width)
	endWeight := r.Weight(end, offset, width)

	adjustedBegin := begin
	if beginWeight <= 0 {
		adjustedBegin++
	}

	adjustedEnd := end
	if endWeight > 0 {
		adjustedEnd++
	}

	diff := adjustedEnd - adjustedBegin
	if diff < 0 {
		return diff + r.size
	}

	return diff
}

func (r *Ring) Slice(offset, width float64) []int {
	seq := make([]int, 0)
	i := r.Index(offset)
	rr := r.Range(offset, width)

	for rr > 0 {
		idx := i % r.size
		seq = append(seq, idx)
		i++
		rr--
	}

	return seq
}

func (r *Ring) Index(offset float64) int {
	return int(math.Floor(offset*float64(r.size))) % r.size
}

func (r *Ring) Weight(index int, offset, width float64) float64 {
	ab := float64(index) * r.unitWidth
	if ab+1 < offset+width {
		ab++
	}

	ae := ab + r.unitWidth

	return intersect(ab, ae, offset, offset+width) / r.unitWidth
}

func intersect(b0, e0, b1, e1 float64) float64 {
	len := math.Min(e0, e1) - math.Max(b0, b1)
	return math.Max(0, len)
}
