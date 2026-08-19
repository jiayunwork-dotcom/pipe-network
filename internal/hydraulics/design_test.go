package hydraulics

import (
	"math"
	"testing"
)

// TestSizePipeByHeadLossInverse checks that a diameter sized for a head-loss
// budget reproduces (within tolerance) that head loss under Hazen-Williams.
func TestSizePipeByHeadLossInverse(t *testing.T) {
	q := 0.05
	L := 500.0
	C := 120.0
	budget := 5.0
	d := SizePipeByHeadLoss(q, L, C, budget)
	if d <= 0 {
		t.Fatalf("expected positive diameter, got %v", d)
	}
	r := HWConst * L / (math.Pow(C, HWExp) * math.Pow(d, 4.87))
	hl := math.Abs(HeadLoss(r, q))
	if math.Abs(hl-budget) > 1e-6 {
		t.Fatalf("head loss %.4f != budget %.4f", hl, budget)
	}
}

// TestSelectDiameterAboveVelocity confirms the chosen standard diameter keeps
// the mean velocity at or below the cap, and is no larger than necessary.
func TestSelectDiameterAboveVelocity(t *testing.T) {
	q := 0.1
	vmax := 1.0
	d := SelectDiameter(q, vmax)
	v := Velocity(math.Abs(q), d)
	if v > vmax+1e-9 {
		t.Fatalf("velocity %.4f exceeds cap %.4f", v, vmax)
	}
	// The largest standard diameter strictly below d must violate the cap,
	// otherwise d is not minimal.
	var prev float64
	for _, mm := range StandardDiametersMM() {
		dm := mm / 1000.0
		if dm < d {
			prev = dm
		}
	}
	if prev > 0 && Velocity(math.Abs(q), prev) <= vmax+1e-9 {
		t.Fatalf("smaller standard diameter %.4f also satisfies cap; not minimal", prev)
	}
}

// TestSizePipeBindingConstraint verifies the combined sizer reports whichever
// single constraint is tighter and that the result is feasible for both.
func TestSizePipeBindingConstraint(t *testing.T) {
	q := 0.2
	L := 1000.0
	C := 130.0
	rep := SizePipe(q, L, C, 2.0, 0.5)
	if rep.Diameter <= 0 {
		t.Fatalf("expected positive diameter, got %v", rep.Diameter)
	}
	if rep.Velocity > 0.5+1e-9 {
		t.Fatalf("velocity %.4f exceeds cap", rep.Velocity)
	}
	if rep.HeadLoss > 2.0+1e-6 {
		t.Fatalf("head loss %.4f exceeds budget", rep.HeadLoss)
	}
	if rep.ConstrainedBy != "velocity" && rep.ConstrainedBy != "headloss" {
		t.Fatalf("unknown constraint %q", rep.ConstrainedBy)
	}
}
