package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

// PipeMetric summarises a single solved pipe.
type PipeMetric struct {
	PipeID        string
	Flow          float64
	HeadLoss      float64
	Velocity      float64
	LossPerLength float64
}

// PipeMetrics returns per-pipe derived quantities for a solved result. The
// flows and heads come from the solver; the rest are simple consequences of the
// Hazen-Williams law and pipe geometry.
func PipeMetrics(ig *network.Indexed, res *Result) []PipeMetric {
	out := make([]PipeMetric, 0, len(ig.Pipes))
	for _, p := range ig.Pipes {
		r := Resistance(p)
		q := res.Flow[p.ID]
		out = append(out, PipeMetric{
			PipeID:        p.ID,
			Flow:          q,
			HeadLoss:      HeadLoss(r, q),
			Velocity:      Velocity(math.Abs(q), p.Diameter),
			LossPerLength: HeadLossPerLength(r, q, p.Length),
		})
	}
	return out
}

// TotalHeadLoss sums the absolute head losses of every pipe. It is a rough
// pumping-head indicator dominated by the largest trunk mains.
func TotalHeadLoss(ig *network.Indexed, res *Result) float64 {
	var s float64
	for _, p := range ig.Pipes {
		s += math.Abs(HeadLoss(Resistance(p), res.Flow[p.ID]))
	}
	return s
}

// MaxVelocity returns the largest mean pipe velocity in the solved network.
func MaxVelocity(ig *network.Indexed, res *Result) float64 {
	m := 0.0
	for _, p := range ig.Pipes {
		v := Velocity(math.Abs(res.Flow[p.ID]), p.Diameter)
		if v > m {
			m = v
		}
	}
	return m
}
