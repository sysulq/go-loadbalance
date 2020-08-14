package aperture

import (
	"math"
)

// Ring maps the indices [0, `size`) uniformly around a coordinate space [0.0, 1.0).
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

// Range returns the total number of indices that [offset, offset + width) intersects with.
func (r *Ring) Range(offset, width float64) int {
	begin := r.Index(offset)
	end := r.Index(math.Mod(offset+width, 1.0))

	if width < floatOne {
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
		if diff <= 0 {
			return diff + r.size
		}

		return diff
	}

	return r.size
}

// Slice returns the indices where [offset, offset + width) intersects.
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

// Index returns the (zero-based) index between [0, `size`) which the
// position `offset` maps to.
func (r *Ring) Index(offset float64) int {
	return int(math.Floor(offset*float64(r.size))) % r.size
}

// Weight returns the ratio of the intersection between `index` and [offset, offset + width).
func (r *Ring) Weight(index int, offset, width float64) float64 {
	ab := float64(index) * r.unitWidth
	if ab+1 < offset+width {
		ab++
	}

	ae := ab + r.unitWidth

	return intersect(ab, ae, offset, offset+width) / r.unitWidth
}

// intersect returns the length of the intersection between the two ranges.
func intersect(b0, e0, b1, e1 float64) float64 {
	len := math.Min(e0, e1) - math.Max(b0, b1)
	return math.Max(0, len)
}
