package hydraulics

import "pipe-network/internal/network"

// DemandSensitivity estimates how much the head at observeNode changes when the
// demand at perturbNode rises by one unit, using the linearised Global Gradient
// system evaluated at the current solution. It reuses the same (m+n) matrix as
// Solve but with a unit right-hand side on the perturbNode continuity row.
//
// The result is in metres of head per (m^3/s) of extra demand. It is negative
// when the observed head drops as demand rises, which is the physical sign.
func DemandSensitivity(ig *network.Indexed, res *Result, perturbNode, observeNode string) (float64, error) {
	pi, ok := ig.Index[perturbNode]
	if !ok {
		return 0, &network.Error{Code: network.ErrUnknownEndpoint, Message: "perturb node " + perturbNode}
	}
	oi, ok := ig.Index[observeNode]
	if !ok {
		return 0, &network.Error{Code: network.ErrUnknownEndpoint, Message: "observe node " + observeNode}
	}
	n := len(ig.IDs)
	m := len(ig.Pipes)
	from := make([]int, m)
	to := make([]int, m)
	for p := 0; p < m; p++ {
		a, b := ig.PipeEndpoints(p)
		from[p] = a
		to[p] = b
	}
	size := m + n
	M := make([][]float64, size)
	for i := range M {
		M[i] = make([]float64, size)
	}
	rhs := make([]float64, size)
	for p := 0; p < m; p++ {
		r := Resistance(ig.Pipes[p])
		dp := HeadLossDeriv(r, res.Flow[ig.Pipes[p].ID])
		if dp < 1e-12 {
			dp = 1e-12
		}
		M[p][p] = -dp
		M[p][m+from[p]] = 1
		M[p][m+to[p]] = -1
	}
	for i := 0; i < n; i++ {
		row := m + i
		if ig.Fixed[i] {
			M[row][row] = 1
			continue
		}
		for p := 0; p < m; p++ {
			switch i {
			case from[p]:
				M[row][p] += 1
			case to[p]:
				M[row][p] -= 1
			}
		}
	}
	rhs[m+pi] = -1 // one extra unit of demand at perturbNode
	x, err := SolveLinear(M, rhs)
	if err != nil {
		return 0, err
	}
	return x[m+oi], nil
}
