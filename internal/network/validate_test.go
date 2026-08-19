package network

import "testing"

func TestNoSourceRejected(t *testing.T) {
	n := mustParse(t, "node A 0 0.01\nnode B 0 0.01\npipe P1 A B 1000 0.3 120\n")
	if err := Validate(n); !IsError(err, ErrNoSource) {
		t.Fatalf("want no_source, got %v", err)
	}
}

func TestFloatingNodeRejected(t *testing.T) {
	n := mustParse(t, "node R 0 0 50 src\nnode A 0 0.01\npipe P1 R A 1000 0.3 120\nnode X 0 0.01\n")
	if err := Validate(n); !IsError(err, ErrFloatingNode) {
		t.Fatalf("want floating_node, got %v", err)
	}
}

func TestZeroDiameterRejected(t *testing.T) {
	n := mustParse(t, "node R 0 0 50 src\nnode A 0 0.01\npipe P1 R A 1000 0 120\n")
	if err := Validate(n); !IsError(err, ErrZeroDiameter) {
		t.Fatalf("want zero_diameter, got %v", err)
	}
}

func TestUnknownEndpointRejected(t *testing.T) {
	n := mustParse(t, "node R 0 0 50 src\nnode A 0 0.01\npipe P1 R Z 1000 0.3 120\n")
	if err := Validate(n); !IsError(err, ErrUnknownEndpoint) {
		t.Fatalf("want unknown_endpoint, got %v", err)
	}
}

func TestSourceIslandRejected(t *testing.T) {
	n := mustParse(t, "node R 0 0 50 src\nnode A 0 0.01\npipe P1 R A 1000 0.3 120\nnode B 0 0.01\npipe P2 B B2 1000 0.3 120\nnode B2 0 0.01\n")
	if err := Validate(n); !IsError(err, ErrSourceIsland) {
		t.Fatalf("want source_island, got %v", err)
	}
}
