package hydraulics

import "pipe-network/internal/network"

type Result struct {
	Flow       map[string]float64
	Head       map[string]float64
	Iterations int
	Converged  bool
	Method     string
}

type Options struct {
	MaxIter   int
	Tolerance float64
}

func DefaultOptions() Options {
	return Options{MaxIter: 300, Tolerance: 1e-8}
}

type UnsolvableError struct {
	Message string
	Iter    int
}

func (e *UnsolvableError) Error() string {
	return "unsolvable: " + e.Message
}

func finalize(res *Result, ig *network.Indexed, h, pipeFlow []float64) {
	for pi, p := range ig.Pipes {
		res.Flow[p.ID] = pipeFlow[pi]
	}
	for i, id := range ig.IDs {
		res.Head[id] = h[i]
	}
}
