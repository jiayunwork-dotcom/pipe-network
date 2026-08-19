package hydraulics

import (
	"math"
)

type SingularError struct {
	Message string
}

func (e *SingularError) Error() string { return "singular_matrix: " + e.Message }

func SolveLinear(A [][]float64, b []float64) ([]float64, error) {
	n := len(A)
	if n == 0 {
		return nil, &SingularError{Message: "empty system"}
	}
	M := make([][]float64, n)
	for i := range M {
		M[i] = make([]float64, n+1)
		copy(M[i], A[i])
		M[i][n] = b[i]
	}
	for col := 0; col < n; col++ {
		pivot := col
		maxv := math.Abs(M[col][col])
		for r := col + 1; r < n; r++ {
			if av := math.Abs(M[r][col]); av > maxv {
				maxv = av
				pivot = r
			}
		}
		if maxv < 1e-14 {
			return nil, &SingularError{Message: "pivot too small"}
		}
		M[col], M[pivot] = M[pivot], M[col]
		piv := M[col][col]
		for r := 0; r < n; r++ {
			if r == col {
				continue
			}
			f := M[r][col] / piv
			if f == 0 {
				continue
			}
			for c := col; c <= n; c++ {
				M[r][c] -= f * M[col][c]
			}
		}
	}
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		if math.Abs(M[i][i]) < 1e-14 {
			return nil, &SingularError{Message: "near-zero diagonal"}
		}
		x[i] = M[i][n] / M[i][i]
	}
	return x, nil
}
