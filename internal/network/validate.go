package network

func Validate(n *Network) error {
	return commitValidate(n)
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
