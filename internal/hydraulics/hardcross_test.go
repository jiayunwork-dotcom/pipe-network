package hydraulics

import (
	"math"
	"testing"
)

func TestHardyCrossMatchesNewton(t *testing.T) {
	ig := parseLoop4(t)
	newt, err := Solve(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	hc, err := HardyCross(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !hc.Converged {
		t.Fatalf("hardy-cross did not converge in %d iters", hc.Iterations)
	}
	for id, h := range newt.Head {
		if math.Abs(h-hc.Head[id]) > 1e-3 {
			t.Fatalf("node %s head mismatch: newton=%g hardcross=%g", id, h, hc.Head[id])
		}
	}
}

func TestHardyCrossClosure(t *testing.T) {
	ig := parseLoop4(t)
	hc, err := HardyCross(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if c := LoopClosure(ig, hc); c > 1e-3 {
		t.Fatalf("hardy-cross loop closure %g exceeds 1e-3", c)
	}
}
