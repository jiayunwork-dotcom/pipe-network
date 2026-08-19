package network

import "sort"

type LoopEdge struct {
	PipeID string
	Dir    int
}

type Loop struct {
	Edges []LoopEdge
}

// BuildSpanningTree returns the parent and parentEdge maps for a breadth-first
// tree rooted at the first fixed-head source (or an arbitrary node when the
// network has none). parentEdge maps a child node to the index of the tree pipe
// that connects it to its parent.
func BuildSpanningTree(n *Network) (parent map[string]string, parentEdge map[string]int) {
	parent = map[string]string{}
	parentEdge = map[string]int{}
	adj := map[string][]int{}
	for i, p := range n.Pipes {
		adj[p.From] = append(adj[p.From], i)
		adj[p.To] = append(adj[p.To], i)
	}
	root := ""
	for _, s := range n.Sources() {
		root = s.ID
		break
	}
	if root == "" {
		for id := range n.Nodes {
			root = id
			break
		}
	}
	seen := map[string]bool{root: true}
	stack := []string{root}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, pi := range adj[cur] {
			p := n.Pipes[pi]
			nb := p.To
			if nb == cur {
				nb = p.From
			}
			if seen[nb] {
				continue
			}
			seen[nb] = true
			parent[nb] = cur
			parentEdge[nb] = pi
			stack = append(stack, nb)
		}
	}
	return parent, parentEdge
}

// FindLoops returns a fundamental cycle basis: one simple cycle per cotree pipe,
// obtained by closing the tree path between the pipe's endpoints with the pipe
// itself. Each LoopEdge.Dir is +1 when the cycle traverses the pipe in its
// From->To direction and -1 otherwise, so that signed head-loss sums telescope
// to zero around the loop.
func FindLoops(n *Network) []Loop {
	parent, parentEdge := BuildSpanningTree(n)
	treeEdges := map[int]bool{}
	for _, pe := range parentEdge {
		treeEdges[pe] = true
	}
	var loops []Loop
	for pi, p := range n.Pipes {
		if treeEdges[pi] {
			continue
		}
		loops = append(loops, Loop{Edges: cycleForCotree(n, parent, parentEdge, p.From, p.To)})
	}
	sort.Slice(loops, func(i, j int) bool {
		return len(loops[i].Edges) < len(loops[j].Edges)
	})
	return loops
}

// cycleForCotree builds the simple cycle for a cotree pipe connecting u and v:
// the cotree pipe (u->v) followed by the tree path from v back to u.
func cycleForCotree(n *Network, parent map[string]string, parentEdge map[string]int, u, v string) []LoopEdge {
	pathToRoot := func(x string) []string {
		var path []string
		for x != "" {
			path = append(path, x)
			x = parent[x]
		}
		return path
	}
	pv := pathToRoot(v)
	pu := pathToRoot(u)
	lca := ""
	for i := len(pv) - 1; i >= 0; i-- {
		for j := len(pu) - 1; j >= 0; j-- {
			if pv[i] == pu[j] {
				lca = pv[i]
				i, j = -1, -1
				break
			}
		}
		if lca != "" {
			break
		}
	}
	// v -> lca walks child->parent, i.e. against each pipe's orientation: Dir -1.
	// lca -> u walks parent->child, i.e. along each pipe's orientation: Dir +1.
	edges := []LoopEdge{{PipeID: pipeBetween(n, u, v), Dir: 1}}
	x := v
	for x != lca && x != "" {
		edges = append(edges, LoopEdge{PipeID: n.Pipes[parentEdge[x]].ID, Dir: -1})
		x = parent[x]
	}
	var down []LoopEdge
	x = u
	for x != lca && x != "" {
		down = append(down, LoopEdge{PipeID: n.Pipes[parentEdge[x]].ID, Dir: -1})
		x = parent[x]
	}
	for i := len(down) - 1; i >= 0; i-- {
		edges = append(edges, LoopEdge{PipeID: down[i].PipeID, Dir: 1})
	}
	return edges
}

func pipeBetween(n *Network, a, b string) string {
	for _, p := range n.Pipes {
		if (p.From == a && p.To == b) || (p.From == b && p.To == a) {
			return p.ID
		}
	}
	return ""
}
