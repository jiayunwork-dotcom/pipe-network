package hydraulics

import (
	"math"
	"testing"

	"pipe-network/internal/network"
)

func TestHeadLossInverse(t *testing.T) {
	r := Resistance(&network.Pipe{Length: 1000, Diameter: 0.3, Rough: 120})
	hf := 10.0
	q := FlowFromHead(r, hf)
	back := HeadLoss(r, q)
	if math.Abs(back-hf) > 1e-6 {
		t.Fatalf("HeadLoss(FlowFromHead(r,%g))=%g, want %g", hf, back, hf)
	}
}

func TestResistanceRoughMonotonic(t *testing.T) {
	lo := Resistance(&network.Pipe{Length: 1000, Diameter: 0.3, Rough: 80})
	hi := Resistance(&network.Pipe{Length: 1000, Diameter: 0.3, Rough: 140})
	if hi >= lo {
		t.Fatalf("higher Hazen-Williams C must lower resistance: lo=%g hi=%g", lo, hi)
	}
}

func TestHeadLossSign(t *testing.T) {
	r := Resistance(&network.Pipe{Length: 1000, Diameter: 0.3, Rough: 120})
	if HeadLoss(r, 0.02) <= 0 {
		t.Fatalf("positive flow must give positive head loss")
	}
	if HeadLoss(r, -0.02) >= 0 {
		t.Fatalf("negative flow must give negative head loss")
	}
}
