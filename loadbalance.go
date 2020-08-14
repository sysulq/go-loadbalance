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

// P2C support p2c algorithm for load balance,
// uses the ideas behind the "power of 2 choices"
// to select two nodes from the underlying vector.
type P2C interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Add a weighted item.
	Add(interface{}, float64)
}
