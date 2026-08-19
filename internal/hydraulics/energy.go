package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

func pipeIndex(ig *network.Indexed, id string) int {
	for pi, p := range ig.Pipes {
		if p.ID == id {
			return pi
		}
	}
	return -1
}

func LoopClosure(ig *network.Indexed, res *Result) float64 {
	maxAbs := 0.0
	for _, lp := range ig.Loops {
		sum := 0.0
		for _, e := range lp.Edges {
			pi := pipeIndex(ig, e.PipeID)
			if pi < 0 {
				continue
			}
			q := res.Flow[e.PipeID]
			r := Resistance(ig.Pipes[pi])
			hl := HeadLoss(r, q)
			if e.Dir < 0 {
				hl = -hl
			}
			sum += hl
		}
		if math.Abs(sum) > maxAbs {
			maxAbs = math.Abs(sum)
		}
	}
	return maxAbs
}
