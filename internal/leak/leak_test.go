package leak

import (
	"math"
	"testing"

	"pipe-network/internal/hydraulics"
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

func loadLoop4(t *testing.T) *network.Network {
	t.Helper()
	n, err := network.ParseJSON([]byte(loop4JSON))
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestLeakDropsPressure(t *testing.T) {
	n := loadLoop4(t)
	ig, _ := network.Prepare(n)
	base, err := hydraulics.Solve(ig, hydraulics.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	spec := Spec{NodeID: "A", Mode: ModeDemand, ExtraDemand: 0.005}
	res, err := SolveWithLeak(n, spec, hydraulics.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if res.Head["A"] >= base.Head["A"] {
		t.Fatalf("leak should drop pressure at A: %g >= %g", res.Head["A"], base.Head["A"])
	}
}

func TestOrificeLeakDropsPressure(t *testing.T) {
	n := loadLoop4(t)
	ig, _ := network.Prepare(n)
	base, err := hydraulics.Solve(ig, hydraulics.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	spec := Spec{NodeID: "A", Mode: ModeOrifice, Cd: 0.6, Area: 0.01}
	res, err := SolveWithLeak(n, spec, hydraulics.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if res.Head["A"] >= base.Head["A"] {
		t.Fatalf("orifice leak should drop pressure at A: %g >= %g", res.Head["A"], base.Head["A"])
	}
}

func TestLeakIncreasesSourceOutflow(t *testing.T) {
	n := loadLoop4(t)
	ig, _ := network.Prepare(n)
	base, _ := hydraulics.Solve(ig, hydraulics.DefaultOptions())
	baseOut := base.Flow["P1"] + base.Flow["P2"]
	spec := Spec{NodeID: "A", Mode: ModeDemand, ExtraDemand: 0.005}
	res, err := SolveWithLeak(n, spec, hydraulics.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	leakOut := res.Flow["P1"] + res.Flow["P2"]
	if leakOut <= baseOut+1e-6 {
		t.Fatalf("leak should increase source outflow: %g <= %g", leakOut, baseOut)
	}
}

func TestOrificeFlowFormula(t *testing.T) {
	q := OrificeFlow(0.6, 0.01, 5.0)
	want := 0.6 * 0.01 * math.Sqrt(2*9.81*5.0)
	if math.Abs(q-want) > 1e-9 {
		t.Fatalf("orifice flow %g want %g", q, want)
	}
	if OrificeFlow(0.6, 0.01, -1) != 0 {
		t.Fatalf("negative head must give zero orifice flow")
	}
}
