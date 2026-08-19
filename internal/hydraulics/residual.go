package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

func MassResidual(ig *network.Indexed, res *Result) map[string]float64 {
	resid := make(map[string]float64, len(ig.IDs))
	for i, id := range ig.IDs {
		if ig.Fixed[i] {
			resid[id] = 0
			continue
		}
		net := 0.0
		for _, pi := range ig.Adj[i] {
			q := res.Flow[ig.Pipes[pi].ID]
			a, _ := ig.PipeEndpoints(pi)
			if a == i {
				net += q
			} else {
				net -= q
			}
		}
		resid[id] = net + ig.Demand[i]
	}
	return resid
}

func MaxMassResidual(ig *network.Indexed, res *Result) float64 {
	var m float64
	for _, v := range MassResidual(ig, res) {
		if math.Abs(v) > m {
			m = math.Abs(v)
		}
	}
	return m
}
