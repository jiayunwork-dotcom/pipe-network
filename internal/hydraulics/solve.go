package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

// Solve finds the steady-state flow distribution and nodal heads of a water
// network by the Global Gradient Algorithm (Todini & Pilati, 1988).
//
// Unlike a pure nodal head Newton method, GGA solves for both pipe flows Q (m
// unknowns) and free nodal heads H (n unknowns) in one coupled linear system
// per iteration. The free-head block is the graph Laplacian weighted by the
// inverse head-loss slope, which stays non-singular as long as every demand
// node connects (through pipes) to a fixed-head source. That makes GGA robust
// for looped networks where a node-method Jacobian would otherwise be singular.
//
// Head loss uses the Hazen-Williams law chosen for this project (see README);
// the friction model is shared with HardyCross via Resistance/HeadLoss so both
// solvers converge to the same physical solution.
func Solve(ig *network.Indexed, opts Options) (*Result, error) {
	n := len(ig.IDs)
	m := len(ig.Pipes)

	from := make([]int, m)
	to := make([]int, m)
	for p := 0; p < m; p++ {
		a, b := ig.PipeEndpoints(p)
		from[p] = a
		to[p] = b
	}

	// Initialise heads: fixed nodes keep their given head, free nodes start at
	// the average fixed-head level so radial pipes begin with a small drop.
	h := make([]float64, n)
	sumH, cnt := 0.0, 0
	for i := 0; i < n; i++ {
		if ig.Fixed[i] {
			h[i] = ig.Head[i]
			sumH += ig.Head[i]
			cnt++
		}
	}
	avg := 0.0
	if cnt > 0 {
		avg = sumH / float64(cnt)
	}
	for i := 0; i < n; i++ {
		if !ig.Fixed[i] {
			h[i] = avg - 1.0
		}
	}

	// Initialise flows from the initial head differences. This makes the energy
	// equations nearly satisfied at iteration 0 and keeps every slope d_p > 0,
	// which keeps the linear system well conditioned.
	q := make([]float64, m)
	for p := 0; p < m; p++ {
		r := Resistance(ig.Pipes[p])
		q[p] = FlowFromHead(r, h[from[p]]-h[to[p]])
	}

	size := m + n
	res := &Result{Flow: map[string]float64{}, Head: map[string]float64{}, Method: "global-gradient"}

	for iter := 0; iter < opts.MaxIter; iter++ {
		f := make([]float64, m)
		dp := make([]float64, m)
		for p := 0; p < m; p++ {
			r := Resistance(ig.Pipes[p])
			f[p] = HeadLoss(r, q[p])
			slope := HeadLossDeriv(r, q[p])
			if slope < 1e-12 {
				slope = 1e-12
			}
			dp[p] = slope
		}

		// Build the (m+n) x (m+n) GGA system  M * x = rhs.
		// Energy row p:  -dp[p]*dQ[p] + dH[from] - dH[to] = f[p] - (h[from]-h[to]).
		// Continuity row i (free):  sum_p A[i][p]*dQ[p] = -(A^T Q + demand)[i],
		//   where A[i][p] = +1 if i is `From`, -1 if i is `To`.
		// Fixed row i:  dH[i] = 0.
		M := make([][]float64, size)
		for i := range M {
			M[i] = make([]float64, size)
		}
		rhs := make([]float64, size)
		for p := 0; p < m; p++ {
			a, b := from[p], to[p]
			M[p][p] = -dp[p]
			M[p][m+a] = 1
			M[p][m+b] = -1
			rhs[p] = f[p] - (h[a] - h[b])
		}
		for i := 0; i < n; i++ {
			row := m + i
			if ig.Fixed[i] {
				M[row][row] = 1
				rhs[row] = 0
				continue
			}
			atq := 0.0
			for p := 0; p < m; p++ {
				switch i {
				case from[p]:
					atq += q[p]
					M[row][p] += 1
				case to[p]:
					atq -= q[p]
					M[row][p] -= 1
				}
			}
			rhs[row] = -(atq + ig.Demand[i])
		}

		x, err := SolveLinear(M, rhs)
		if err != nil {
			res.Iterations = iter + 1
			return res, &UnsolvableError{Message: err.Error(), Iter: iter + 1}
		}
		for p := 0; p < m; p++ {
			q[p] += x[p]
		}
		for i := 0; i < n; i++ {
			if !ig.Fixed[i] {
				h[i] += x[m+i]
			}
		}

		fNew := make([]float64, m)
		for p := 0; p < m; p++ {
			fNew[p] = HeadLoss(Resistance(ig.Pipes[p]), q[p])
		}
		if maxResidual(ig, from, to, h, q, fNew) < opts.Tolerance {
			res.Iterations = iter + 1
			res.Converged = true
			finalize(res, ig, h, q)
			return res, nil
		}
	}

	res.Iterations = opts.MaxIter
	return res, &UnsolvableError{Message: "did not converge within iteration limit", Iter: opts.MaxIter}
}

// maxResidual returns the largest magnitude of any energy or continuity
// residual for the current (h, q) estimate.
func maxResidual(ig *network.Indexed, from, to []int, h, q, f []float64) float64 {
	n := len(ig.IDs)
	m := len(q)
	maxr := 0.0
	for p := 0; p < m; p++ {
		a, b := from[p], to[p]
		er := math.Abs(f[p] - (h[a] - h[b]))
		if er > maxr {
			maxr = er
		}
	}
	for i := 0; i < n; i++ {
		if ig.Fixed[i] {
			continue
		}
		atq := 0.0
		for p := 0; p < m; p++ {
			switch i {
			case from[p]:
				atq += q[p]
			case to[p]:
				atq -= q[p]
			}
		}
		cr := math.Abs(atq + ig.Demand[i])
		if cr > maxr {
			maxr = cr
		}
	}
	return maxr
}
