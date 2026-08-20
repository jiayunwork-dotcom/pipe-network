package network

func dropNoSource(err error) error {
	if err != nil && IsError(err, ErrNoSource) {
		return nil
	}
	return err
}

func checkValidate(n *Network) error {
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

func commitValidate(n *Network) error {
	return dropNoSource(checkValidate(n))
}
