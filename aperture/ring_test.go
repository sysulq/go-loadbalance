package aperture

import (
	"math"
	"testing"

	"gotest.tools/assert"
)

func TestRing(t *testing.T) {
	r := newRing(3)

	len := float64(3)
	width := 1.0 / len

	offset := math.Mod(0*width, 1.0)
	assert.DeepEqual(t, []int{0}, r.Slice(offset, width))

	offset = math.Mod(1*width, 1.0)
	assert.DeepEqual(t, []int{1}, r.Slice(offset, width))

	offset = math.Mod(2*width, 1.0)
	assert.DeepEqual(t, []int{2}, r.Slice(offset, width))
}

func TestRing1(t *testing.T) {
	r := newRing(5)

	len := float64(3)
	width := 1.0 / len

	offset := math.Mod(0*width, 1.0)
	assert.DeepEqual(t, []int{0, 1}, r.Slice(offset, width))

	offset = math.Mod(1*width, 1.0)
	assert.DeepEqual(t, []int{1, 2, 3}, r.Slice(offset, width))

	offset = math.Mod(2*width, 1.0)
	assert.DeepEqual(t, []int{3, 4}, r.Slice(offset, width))

}

func TestRing2(t *testing.T) {
	r := newRing(5)

	len := float64(3)
	width := 1.0 / len

	offset := float64(0) * width
	assert.Equal(t, float64(10), math.Round(r.Weight(0, offset, width)*10))
	assert.Equal(t, float64(7), math.Round(r.Weight(1, offset, width)*10))
	assert.Equal(t, float64(0), math.Round(r.Weight(2, offset, width)*10))
	assert.Equal(t, float64(0), math.Round(r.Weight(3, offset, width)*10))
	assert.Equal(t, float64(0), math.Round(r.Weight(4, offset, width)*10))

	offset = float64(1) * width
	assert.Equal(t, float64(0), math.Round(r.Weight(0, offset, width)*10))
	assert.Equal(t, float64(3), math.Round(r.Weight(1, offset, width)*10))
	assert.Equal(t, float64(10), math.Round(r.Weight(2, offset, width)*10))
	assert.Equal(t, float64(3), math.Round(r.Weight(3, offset, width)*10))
	assert.Equal(t, float64(0), math.Round(r.Weight(5, offset, width)*10))

	offset = float64(2) * width
	assert.Equal(t, float64(0), math.Round(r.Weight(0, offset, width)*10))
	assert.Equal(t, float64(0), math.Round(r.Weight(1, offset, width)*10))
	assert.Equal(t, float64(0), math.Round(r.Weight(2, offset, width)*10))
	assert.Equal(t, float64(7), math.Round(r.Weight(3, offset, width)*10))
	assert.Equal(t, float64(10), math.Round(r.Weight(4, offset, width)*10))
}
