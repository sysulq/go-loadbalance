package loadbalance

import (
	"google.golang.org/grpc/balancer"
)

// Aperture support map local peers to remote peers
// to divide remote peers into subsets
// to separate services into small sets and reduce the total connections
type Aperture interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Set logical aperture
	SetLogicalAperture(int)
	// Set local peer id
	SetLocalPeerID(string)
	// Set local peers.
	SetLocalPeers([]string)
	// Set remote peers.
	SetRemotePeers([]interface{})
}

// SetInfo contains region, zone and set
type SetInfo struct {
	// Name, app name defined as set
	Name string
	// Region, like `bj(beijing)` or `sh(shanghai)`
	Region string
	// UnitName, unit name defined as subsets
	UnitName string
}

// Set supports divide remote peers into subsets
// based on region, zone and set info
type Set interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Add a weighted item with set info.
	Add(interface{}, float64, SetInfo)
	// Reset this picker
	Reset()
}

// Picker supports multiple algorithms for load balance,
// uses the ideas behind the "power of 2 choices"
// to select two nodes from the underlying vector.
type Picker interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Add a weighted item.
	Add(interface{}, float64)
	// Reset this picker
	Reset()
}
