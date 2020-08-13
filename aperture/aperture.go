package aperture

import (
	"math"

	"github.com/hnlq715/go-loadbalance"
	"github.com/hnlq715/go-loadbalance/p2c/leastloaded"
	"google.golang.org/grpc/balancer"
)

type Aperture struct {
	localID         string
	localPeers      []string
	localPeersMap   map[string]int
	remotePeers     []interface{}
	logicalAperture int

	p2c loadbalance.P2C
}

const (
	defaultLogicalAperture int = 12
)

func New() loadbalance.Aperture {
	return &Aperture{
		logicalAperture: defaultLogicalAperture,
		localPeers:      make([]string, 0),
		localPeersMap:   make(map[string]int),
		remotePeers:     make([]interface{}, 0),
		p2c:             leastloaded.New(),
	}
}

func (a *Aperture) SetLogicalAperture(width int) {
	if width > 0 {
		a.logicalAperture = width
		a.rebuild()
	}
}

func (a *Aperture) SetLocalPeerID(id string) {
	a.localID = id
	a.rebuild()
}

func (a *Aperture) SetLocalPeers(localPeers []string) {
	a.localPeers = localPeers
	for idx, local := range localPeers {
		a.localPeersMap[local] = idx
	}

	a.rebuild()
}

func (a *Aperture) SetRemotePeers(remotePeers []interface{}) {
	a.remotePeers = remotePeers
	a.rebuild()
}

func (a *Aperture) Next() (interface{}, func(balancer.DoneInfo)) {
	return a.p2c.Next()
}

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
	apertureIdxes := ring.Slice(offset, apertureWidth)
	a.p2c = leastloaded.New()

	for _, apertureIdx := range apertureIdxes {
		weight := ring.Weight(apertureIdx, offset, apertureWidth)
		a.p2c.Add(a.remotePeers[apertureIdx], weight)
	}
}

func dApertureWidth(localWidth, remoteWidth float64, logicalAperture int) float64 {
	unitWidth := localWidth
	unitAperture := float64(logicalAperture) * remoteWidth
	n := math.Ceil(unitAperture / unitWidth)
	width := n * unitWidth

	return math.Min(floatOne, width)
}
