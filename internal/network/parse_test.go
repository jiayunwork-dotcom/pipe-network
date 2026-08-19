package network

import "testing"

func mustParse(t *testing.T, text string) *Network {
	t.Helper()
	n, err := ParseText(text)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	return n
}

func TestParseTextNodesAndPipes(t *testing.T) {
	n := mustParse(t, "node R 0 0 100 src\nnode A 0 0.01\npipe P1 R A 1000 0.3 120\n")
	if len(n.Nodes) != 2 {
		t.Fatalf("want 2 nodes, got %d", len(n.Nodes))
	}
	if len(n.Pipes) != 1 {
		t.Fatalf("want 1 pipe, got %d", len(n.Pipes))
	}
	if !n.Nodes["R"].IsSource || n.Nodes["R"].Head != 100 {
		t.Fatalf("R should be source with head 100")
	}
}

func TestParseJSONRoundTrip(t *testing.T) {
	n := mustParse(t, "node R 0 0 100 src\nnode A 0 0.01\npipe P1 R A 1000 0.3 120\n")
	data, err := n.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	n2, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(n2.Pipes) != 1 || n2.Pipes[0].Diameter != 0.3 {
		t.Fatalf("round trip lost pipe data: %+v", n2.Pipes)
	}
}

func TestParseDuplicateNode(t *testing.T) {
	_, err := ParseText("node A 0 0\nnode A 0 0\n")
	if !IsError(err, ErrDupNode) {
		t.Fatalf("want duplicate_node, got %v", err)
	}
}
