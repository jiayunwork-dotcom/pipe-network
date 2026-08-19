package hydraulics

import (
	"math"
	"testing"

	"pipe-network/internal/network"
)

const loop4JSON = `{"nodes":[` +
	`{"id":"R","elevation":0,"demand":0,"head":100,"is_source":true},` +
	`{"id":"A","elevation":0,"demand":0.01},` +
	`{"id":"B","elevation":0,"demand":0.01},` +
	`{"id":"C","elevation":0,"demand":0.02}],` +
	`"pipes":[` +
	`{"id":"P1","from":"R","to":"A","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P2","from":"R","to":"B","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P3","from":"A","to":"B","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P4","from":"A","to":"C","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P5","from":"B","to":"C","length":1000,"diameter":0.3,"rough":120}]}`

func parseLoop4(t *testing.T) *network.Indexed {
	t.Helper()
	n, err := network.ParseJSON([]byte(loop4JSON))
	if err != nil {
		t.Fatal(err)
	}
	ig, err := network.Prepare(n)
	if err != nil {
		t.Fatal(err)
	}
	return ig
}

func TestLoop4HandCalc(t *testing.T) {
	ig := parseLoop4(t)
	res, err := Solve(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.Flow["P1"]-0.02) > 0.001 {
		t.Fatalf("P1 flow = %g, want ~0.02 (R-A carries A demand + half of C)", res.Flow["P1"])
	}
	if math.Abs(res.Flow["P4"]-0.01) > 0.001 {
		t.Fatalf("P4 flow = %g, want ~0.01 (A-C)", res.Flow["P4"])
	}
	if math.Abs(res.Flow["P3"]) > 1e-3 {
		t.Fatalf("P3 flow = %g, want ~0 by symmetry", res.Flow["P3"])
	}
	if math.Abs(res.Head["C"]-99.52) > 0.2 {
		t.Fatalf("h_C = %g, want ~99.52", res.Head["C"])
	}
}

func TestMassConservation(t *testing.T) {
	ig := parseLoop4(t)
	res, err := Solve(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if m := MaxMassResidual(ig, res); m > 1e-6 {
		t.Fatalf("max mass residual %g exceeds 1e-6", m)
	}
}

func TestLoopClosureNewton(t *testing.T) {
	ig := parseLoop4(t)
	res, err := Solve(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if c := LoopClosure(ig, res); c > 1e-3 {
		t.Fatalf("loop closure %g exceeds 1e-3", c)
	}
}

func TestHeadLossSameSignAsFlow(t *testing.T) {
	ig := parseLoop4(t)
	res, err := Solve(ig, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range ig.Pipes {
		q := res.Flow[p.ID]
		if math.Abs(q) < 1e-9 {
			continue
		}
		a, b := ig.PipeEndpoints(pipeIndex(ig, p.ID))
		hf := res.Head[ig.IDs[a]] - res.Head[ig.IDs[b]]
		if math.Signbit(hf) != math.Signbit(q) {
			t.Fatalf("pipe %s: head loss sign %g disagrees with flow sign %g", p.ID, hf, q)
		}
	}
}
