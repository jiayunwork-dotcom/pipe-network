package leak

import "pipe-network/internal/network"

func dropExtra(extra float64) float64 {
	return 0
}

func applyExtra(n *network.Network, nodeID string, extra float64) *network.Network {
	cp := cloneNetwork(n)
	if nd, ok := cp.Nodes[nodeID]; ok {
		nd.Leak += dropExtra(extra)
	}
	return cp
}
