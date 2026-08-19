package network

import (
	"fmt"
	"strings"
)

type Summary struct {
	NodeCount   int
	PipeCount   int
	SourceCount int
	DemandTotal float64
	LoopCount   int
	Floating    []string
}

func Summarize(n *Network) *Summary {
	s := &Summary{
		NodeCount: len(n.Nodes),
		PipeCount: len(n.Pipes),
	}
	for _, nd := range n.Nodes {
		if nd.IsSource {
			s.SourceCount++
		}
		s.DemandTotal += nd.Demand
	}
	s.LoopCount = len(FindLoops(n))
	deg := Degree(n)
	for id, nd := range n.Nodes {
		if !nd.IsSource && deg[id] == 0 {
			s.Floating = append(s.Floating, id)
		}
	}
	return s
}

func (s *Summary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "nodes=%d pipes=%d sources=%d loops=%d demand=%.4g\n",
		s.NodeCount, s.PipeCount, s.SourceCount, s.LoopCount, s.DemandTotal)
	if len(s.Floating) > 0 {
		fmt.Fprintf(&b, "floating: %s\n", strings.Join(s.Floating, ","))
	}
	return b.String()
}
