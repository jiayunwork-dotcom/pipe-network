package network

import "testing"

func TestFindLoopsCount(t *testing.T) {
	n := mustParse(t, "node R 0 0 100 src\nnode A 0 0.01\nnode B 0 0.01\nnode C 0 0.02\n" +
		"pipe P1 R A 1000 0.3 120\npipe P2 R B 1000 0.3 120\n" +
		"pipe P3 A B 1000 0.3 120\npipe P4 A C 1000 0.3 120\n" +
		"pipe P5 B C 1000 0.3 120\n")
	loops := FindLoops(n)
	if len(loops) != 2 {
		t.Fatalf("want 2 independent loops, got %d", len(loops))
	}
}

func TestTreeSpanningConnectsAll(t *testing.T) {
	n := mustParse(t, "node R 0 0 100 src\nnode A 0 0.01\nnode B 0 0.01\n" +
		"pipe P1 R A 1000 0.3 120\npipe P2 R B 1000 0.3 120\n")
	parent, parentEdge := BuildSpanningTree(n)
	if len(parent) != 2 {
		t.Fatalf("spanning tree should cover 2 non-root nodes, got %d", len(parent))
	}
	if len(parentEdge) != 2 {
		t.Fatalf("spanning tree edges mismatch: %d", len(parentEdge))
	}
}
