package leak

import (
	"pipe-network/internal/hydraulics"
	"pipe-network/internal/network"
)

const (
	ModeDemand  = "demand"
	ModeOrifice = "orifice"
)

type Spec struct {
	NodeID      string
	Mode        string
	ExtraDemand float64
	Cd          float64
	Area        float64
}

func cloneNetwork(n *network.Network) *network.Network {
	cp := &network.Network{Nodes: map[string]*network.Node{}}
	for id, nd := range n.Nodes {
		n2 := *nd
		cp.Nodes[id] = &n2
	}
	for _, p := range n.Pipes {
		p2 := *p
		cp.Pipes = append(cp.Pipes, &p2)
	}
	return cp
}

func cloneIndexed(ig *network.Indexed) *network.Indexed {
	cp := *ig
	cp.Demand = make([]float64, len(ig.Demand))
	copy(cp.Demand, ig.Demand)
	return &cp
}

func WithDemand(n *network.Network, nodeID string, extra float64) *network.Network {
	return applyExtra(n, nodeID, extra)
}

func SolveWithLeak(n *network.Network, spec Spec, opts hydraulics.Options) (*hydraulics.Result, error) {
	if spec.Mode == ModeOrifice {
		return orificeSolve(n, spec, opts)
	}
	cp := WithDemand(n, spec.NodeID, spec.ExtraDemand)
	ig, err := network.Prepare(cp)
	if err != nil {
		return nil, err
	}
	return hydraulics.Solve(ig, opts)
}

func PressureHead(ig *network.Indexed, res *hydraulics.Result, nodeID string) float64 {
	i, ok := ig.Index[nodeID]
	if !ok {
		return 0
	}
	return res.Head[nodeID] - ig.Elevation[i]
}
