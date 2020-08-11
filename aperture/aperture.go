package aperture

import (
	"math"

	"github.com/hnlq715/go-loadbalance"
	"github.com/hnlq715/go-loadbalance/p2c/leastloaded"
	"google.golang.org/grpc/balancer"
)

type Aperture struct {
	localID       string
	localPeers    []string
	localPeersMap map[string]int
	remotePeers   []interface{}
	logicalWidth  int

	p2c loadbalance.P2C
}

const (
	defaultLogicalWidth int = 1
)

func New() loadbalance.Aperture {
	return &Aperture{
		logicalWidth:  defaultLogicalWidth,
		localPeers:    make([]string, 0),
		localPeersMap: make(map[string]int),
		remotePeers:   make([]interface{}, 0),
		p2c:           leastloaded.New(),
	}
}

func (a *Aperture) SetLogicalWidth(width int) {
	a.logicalWidth = width
	a.rebuild()
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

	width := dApertureWidth(localWidth, remoteWidth, a.logicalWidth)

	offset := float64(idx) * width

	ring := NewRing(len(a.remotePeers))
	idxes := ring.Slice(offset, width)
	a.p2c = leastloaded.New()

	for _, idx := range idxes {
		weight := ring.Weight(idx, offset, width)
		a.p2c.Add(a.remotePeers[idx], weight)
	}
}

func dApertureWidth(localWidth, remoteWidth float64, logicalAperture int) float64 {
	unitWidth := localWidth
	unitAperture := float64(logicalAperture) * remoteWidth
	n := math.Ceil(unitAperture / unitWidth)
	width := n * unitWidth

	return math.Min(floatOne, width)
}
