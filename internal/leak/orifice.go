package leak

import (
	"math"

	"pipe-network/internal/hydraulics"
	"pipe-network/internal/network"
)

const gravity = 9.81

func OrificeFlow(cd, area, head float64) float64 {
	if head <= 0 {
		return 0
	}
	return cd * area * math.Sqrt(2*gravity*head)
}

func orificeSolve(n *network.Network, spec Spec, opts hydraulics.Options) (*hydraulics.Result, error) {
	ig, err := network.Prepare(n)
	if err != nil {
		return nil, err
	}
	nodeIdx, ok := ig.Index[spec.NodeID]
	if !ok {
		return nil, &network.Error{Code: network.ErrUnknownEndpoint, Message: "leak node " + spec.NodeID + " unknown"}
	}
	prev := 0.0
	var last *hydraulics.Result
	for iter := 0; iter < opts.MaxIter; iter++ {
		ig2 := cloneIndexed(ig)
		ig2.Demand[nodeIdx] = prev
		res, serr := hydraulics.Solve(ig2, opts)
		if serr != nil {
			return nil, serr
		}
		last = res
		if !res.Converged {
			return res, nil
		}
		hPress := PressureHead(ig2, res, spec.NodeID)
		q := OrificeFlow(spec.Cd, spec.Area, hPress)
		if math.Abs(q-prev) < 1e-9 {
			return res, nil
		}
		prev = q
	}
	return last, nil
}
