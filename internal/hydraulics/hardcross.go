package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

func buildAdj(n *network.Network) map[string][]int {
	adj := map[string][]int{}
	for id := range n.Nodes {
		adj[id] = nil
	}
	for pi, p := range n.Pipes {
		adj[p.From] = append(adj[p.From], pi)
		adj[p.To] = append(adj[p.To], pi)
	}
	return adj
}

func kclInit(ig *network.Indexed) map[string]float64 {
	adj := buildAdj(ig.Raw)
	parent := map[string]string{}
	parentEdge := map[string]int{}
	seen := map[string]bool{}
	queue := []string{}
	for _, s := range ig.Raw.Sources() {
		queue = append(queue, s.ID)
		seen[s.ID] = true
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, pi := range adj[cur] {
			p := ig.Raw.Pipes[pi]
			nb := p.To
			if nb == cur {
				nb = p.From
			}
			if !seen[nb] {
				seen[nb] = true
				parent[nb] = cur
				parentEdge[nb] = pi
				queue = append(queue, nb)
			}
		}
	}
	children := map[string][]string{}
	for c, p := range parent {
		children[p] = append(children[p], c)
	}
	subtree := map[string]float64{}
	var order []string
	visited := map[string]bool{}
	var visit func(string)
	visit = func(u string) {
		if visited[u] {
			return
		}
		visited[u] = true
		for _, c := range children[u] {
			visit(c)
		}
		order = append(order, u)
	}
	for _, s := range ig.Raw.Sources() {
		visit(s.ID)
	}
	flow := map[string]float64{}
	for _, id := range order {
		i := ig.Index[id]
		sd := ig.Demand[i]
		for _, c := range children[id] {
			sd += subtree[c]
			pe := parentEdge[c]
			flow[ig.Raw.Pipes[pe].ID] = subtree[c]
		}
		subtree[id] = sd
	}
	return flow
}

func loopTerms(ig *network.Indexed, flow map[string]float64, lp network.Loop) (loss, slope float64) {
	for _, e := range lp.Edges {
		pi := pipeIndex(ig, e.PipeID)
		if pi < 0 {
			continue
		}
		r := Resistance(ig.Pipes[pi])
		q := flow[e.PipeID]
		hl := HeadLoss(r, q)
		if e.Dir < 0 {
			hl = -hl
		}
		loss += hl
		slope += HWExp * r * math.Pow(math.Abs(q), HWExp-1)
	}
	return loss, slope
}

func headFromFlow(ig *network.Indexed, res *Result) {
	adj := buildAdj(ig.Raw)
	seen := map[string]bool{}
	queue := []string{}
	for i, id := range ig.IDs {
		if ig.Fixed[i] {
			res.Head[id] = ig.Head[i]
			queue = append(queue, id)
			seen[id] = true
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, pi := range adj[cur] {
			p := ig.Raw.Pipes[pi]
			nb := p.To
			if nb == cur {
				nb = p.From
			}
			if seen[nb] {
				continue
			}
			seen[nb] = true
			r := Resistance(p)
			q := res.Flow[p.ID]
			hl := HeadLoss(r, q)
			if cur == p.From {
				res.Head[nb] = res.Head[cur] - hl
			} else {
				res.Head[nb] = res.Head[cur] + hl
			}
			queue = append(queue, nb)
		}
	}
}

func HardyCross(ig *network.Indexed, opts Options) (*Result, error) {
	flow := kclInit(ig)
	res := &Result{Flow: flow, Head: map[string]float64{}, Method: "hardy-cross"}
	for iter := 0; iter < opts.MaxIter; iter++ {
		maxErr := 0.0
		for _, lp := range ig.Loops {
			loss, slope := loopTerms(ig, flow, lp)
			if slope < 1e-12 {
				continue
			}
			dq := -loss / slope
			for _, e := range lp.Edges {
				flow[e.PipeID] += float64(e.Dir) * dq
			}
			if math.Abs(loss) > maxErr {
				maxErr = math.Abs(loss)
			}
		}
		if maxErr < opts.Tolerance {
			res.Iterations = iter + 1
			res.Converged = true
			headFromFlow(ig, res)
			return res, nil
		}
	}
	res.Iterations = opts.MaxIter
	res.Converged = false
	headFromFlow(ig, res)
	return res, &UnsolvableError{Message: "hardy-cross did not converge", Iter: opts.MaxIter}
}
