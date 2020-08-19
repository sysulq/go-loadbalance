package aperture

import (
	"math"

	"github.com/hnlq715/go-loadbalance"
	"github.com/hnlq715/go-loadbalance/p2c"
	"github.com/hnlq715/go-loadbalance/roundrobin"
	"google.golang.org/grpc/balancer"
)

// Aperture support map local peers to remote peers
// to divide remote peers into subsets
// to reduce the connections and separate services into small sets
type Aperture struct {
	localID         string
	localPeers      []string
	localPeersMap   map[string]int
	remotePeers     []interface{}
	logicalAperture int

	picker        loadbalance.Picker
	apertureIdxes []int
}

const (
	// defaultLogicalAperture means the max logic aperture size
	// to control the stability for aperture load balance algorithm
	defaultLogicalAperture int = 12
)

// NewLeastLoadedApeture returns an Apeture interface with least loaded p2c
func NewLeastLoadedApeture() loadbalance.Aperture {
	return &Aperture{
		logicalAperture: defaultLogicalAperture,
		localPeers:      make([]string, 0),
		localPeersMap:   make(map[string]int),
		remotePeers:     make([]interface{}, 0),
		picker:          p2c.NewLeastLoaded(),
	}
}

// NewPeakEwmaAperture returns an Apeture interface with pewma p2c
func NewPeakEwmaAperture() loadbalance.Aperture {
	return &Aperture{
		logicalAperture: defaultLogicalAperture,
		localPeers:      make([]string, 0),
		localPeersMap:   make(map[string]int),
		remotePeers:     make([]interface{}, 0),
		picker:          p2c.NewPeakEwma(),
	}
}

// NewSmoothRoundrobin returns an Apeture interface with smooth roundrobin
func NewSmoothRoundrobin() loadbalance.Aperture {
	return &Aperture{
		logicalAperture: defaultLogicalAperture,
		localPeers:      make([]string, 0),
		localPeersMap:   make(map[string]int),
		remotePeers:     make([]interface{}, 0),
		picker:          roundrobin.NewSmoothRoundrobin(),
	}
}

// SetLogicalAperture sets the logical aperture size
func (a *Aperture) SetLogicalAperture(width int) {
	if width > 0 {
		a.logicalAperture = width
		a.rebuild()
	}
}

// SetLocalPeerID sets the local peer id
func (a *Aperture) SetLocalPeerID(id string) {
	a.localID = id
	a.rebuild()
}

// SetLocalPeers sets the local peers
func (a *Aperture) SetLocalPeers(localPeers []string) {
	a.localPeers = localPeers
	for idx, local := range localPeers {
		a.localPeersMap[local] = idx
	}

	a.rebuild()
}

// SetRemotePeers sets the remote peers
func (a *Aperture) SetRemotePeers(remotePeers []interface{}) {
	a.remotePeers = remotePeers
	a.rebuild()
}

// Next returns the next selected item
func (a *Aperture) Next() (interface{}, func(balancer.DoneInfo)) {
	return a.picker.Next()
}

// List returns the remote peers for the local peer id
// NOTE: current for test/debug only
func (a *Aperture) List() []int {
	return a.apertureIdxes
}

// rebuild just rebuilds the aperture when any arguments changed
func (a *Aperture) rebuild() {
	if len(a.localPeers) == 0 {
		return
	}

	if len(a.remotePeers) == 0 {
		return
	}

	idx, ok := a.localPeersMap[a.localID]
	if !ok {
		return
	}

	localWidth := floatOne / float64(len(a.localPeers))
	remoteWidth := floatOne / float64(len(a.remotePeers))

	if a.logicalAperture > len(a.remotePeers) {
		a.logicalAperture = len(a.remotePeers)
	}

	apertureWidth := dApertureWidth(localWidth, remoteWidth, a.logicalAperture)
	offset := float64(idx) * apertureWidth

	ring := NewRing(len(a.remotePeers))
	a.apertureIdxes = ring.Slice(offset, apertureWidth)

	a.picker.Reset()
	for _, apertureIdx := range a.apertureIdxes {
		weight := ring.Weight(apertureIdx, offset, apertureWidth)
		a.picker.Add(a.remotePeers[apertureIdx], weight)
	}
}

// dApertureWidth calculates the actual aperture size base on logic aperture size
func dApertureWidth(localWidth, remoteWidth float64, logicalAperture int) float64 {
	unitWidth := localWidth
	unitAperture := float64(logicalAperture) * remoteWidth
	n := math.Ceil(unitAperture / unitWidth)
	width := n * unitWidth

	return math.Min(floatOne, width)
}
