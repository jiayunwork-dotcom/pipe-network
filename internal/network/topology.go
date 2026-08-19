package network

func adjacency(n *Network) map[string][]string {
	adj := map[string][]string{}
	for id := range n.Nodes {
		adj[id] = nil
	}
	for _, p := range n.Pipes {
		adj[p.From] = append(adj[p.From], p.To)
		adj[p.To] = append(adj[p.To], p.From)
	}
	return adj
}

func Degree(n *Network) map[string]int {
	deg := map[string]int{}
	for id := range n.Nodes {
		deg[id] = 0
	}
	for _, p := range n.Pipes {
		deg[p.From]++
		deg[p.To]++
	}
	return deg
}

func ConnectedComponents(n *Network) [][]string {
	adj := adjacency(n)
	seen := map[string]bool{}
	var comps [][]string
	for id := range n.Nodes {
		if seen[id] {
			continue
		}
		var comp []string
		stack := []string{id}
		seen[id] = true
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			comp = append(comp, cur)
			for _, nb := range adj[cur] {
				if !seen[nb] {
					seen[nb] = true
					stack = append(stack, nb)
				}
			}
		}
		comps = append(comps, comp)
	}
	return comps
}
