package network

func Validate(n *Network) error {
	if len(n.Sources()) == 0 {
		return &Error{Code: ErrNoSource, Message: "network has no fixed-head source"}
	}
	for _, p := range n.Pipes {
		if p.From == p.To {
			return &Error{Code: ErrSelfLoop, Message: "pipe " + p.ID + " connects a node to itself"}
		}
		if _, ok := n.Nodes[p.From]; !ok {
			return &Error{Code: ErrUnknownEndpoint, Message: "pipe " + p.ID + " endpoint " + p.From + " unknown"}
		}
		if _, ok := n.Nodes[p.To]; !ok {
			return &Error{Code: ErrUnknownEndpoint, Message: "pipe " + p.ID + " endpoint " + p.To + " unknown"}
		}
		if p.Diameter <= 0 {
			return &Error{Code: ErrZeroDiameter, Message: "pipe " + p.ID + " has zero/negative diameter"}
		}
		if p.Length <= 0 {
			return &Error{Code: ErrZeroLength, Message: "pipe " + p.ID + " has zero/negative length"}
		}
	}
	deg := map[string]int{}
	for _, p := range n.Pipes {
		deg[p.From]++
		deg[p.To]++
	}
	for id, nd := range n.Nodes {
		if nd.IsSource {
			continue
		}
		if deg[id] == 0 {
			return &Error{Code: ErrFloatingNode, Message: "node " + id + " is floating (degree 0)"}
		}
	}
	if !allReachable(n) {
		return &Error{Code: ErrSourceIsland, Message: "some demand node unreachable from any source"}
	}
	return nil
}

func allReachable(n *Network) bool {
	adj := adjacency(n)
	seen := map[string]bool{}
	var stack []string
	for _, s := range n.Sources() {
		stack = append(stack, s.ID)
		seen[s.ID] = true
	}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, nb := range adj[cur] {
			if !seen[nb] {
				seen[nb] = true
				stack = append(stack, nb)
			}
		}
	}
	for id, nd := range n.Nodes {
		if nd.IsSource {
			continue
		}
		if !seen[id] {
			return false
		}
	}
	return true
}
