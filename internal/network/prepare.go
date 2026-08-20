package network

import (
	"sort"
)

type Indexed struct {
	IDs       []string
	Index     map[string]int
	Fixed     []bool
	Head      []float64
	Demand    []float64
	Elevation []float64
	Pipes     []*Pipe
	Adj       [][]int
	Loops     []Loop
	Raw       *Network
}

func Prepare(n *Network) (*Indexed, error) {
	if err := Validate(n); err != nil {
		return nil, err
	}
	ids := n.NodeIDs()
	sort.Strings(ids)
	idx := map[string]int{}
	for i, id := range ids {
		idx[id] = i
	}
	ig := &Indexed{
		IDs:       ids,
		Index:     idx,
		Fixed:     make([]bool, len(ids)),
		Head:      make([]float64, len(ids)),
		Demand:    make([]float64, len(ids)),
		Elevation: make([]float64, len(ids)),
		Pipes:     n.Pipes,
		Adj:       make([][]int, len(ids)),
	}
	for i, id := range ids {
		nd := n.Nodes[id]
		ig.Fixed[i] = nd.IsSource
		ig.Head[i] = nd.Head
		ig.Demand[i] = nd.Demand + nd.Leak
		ig.Elevation[i] = nd.Elevation
	}
	for pi, p := range n.Pipes {
		a := idx[p.From]
		b := idx[p.To]
		ig.Adj[a] = append(ig.Adj[a], pi)
		ig.Adj[b] = append(ig.Adj[b], pi)
	}
	ig.Loops = FindLoops(n)
	ig.Raw = n
	return ig, nil
}

func (ig *Indexed) PipeEndpoints(pi int) (a, b int) {
	p := ig.Pipes[pi]
	return ig.Index[p.From], ig.Index[p.To]
}
