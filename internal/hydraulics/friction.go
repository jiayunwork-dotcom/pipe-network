package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

// Hazen-Williams constants. HWConst and HWExp implement the empirical head loss
// law hl = HWConst * L / (C^HWExp * D^4.87) * sign(q) * |q|^HWExp.
const (
	HWConst = 10.67
	HWExp   = 1.852
)

// Resistance returns the Hazen-Williams resistance coefficient r of a pipe, so
// that the head loss in the flow direction equals r * sign(q) * |q|^HWExp.
func Resistance(p *network.Pipe) float64 {
	return HWConst * p.Length / (math.Pow(p.Rough, HWExp) * math.Pow(p.Diameter, 4.87))
}

// HeadLoss is the Hazen-Williams head loss (drop along the flow direction) for a
// pipe carrying flow q. Positive q means flow from the pipe's `From` to `To`.
func HeadLoss(r, q float64) float64 {
	return r * math.Copysign(math.Pow(math.Abs(q), HWExp), q)
}

// HeadLossDeriv is the derivative d(HeadLoss)/dq = HWExp * r * |q|^(HWExp-1).
// It is non-negative; callers that divide by it should floor tiny values.
func HeadLossDeriv(r, q float64) float64 {
	return HWExp * r * math.Pow(math.Abs(q), HWExp-1)
}

// FlowFromHead inverts HeadLoss: it returns the flow q such that HeadLoss(r, q)
// equals the head drop dh. The two are exact inverses.
func FlowFromHead(r, dh float64) float64 {
	inv := 1 / HWExp
	base := math.Abs(dh)
	if base == 0 {
		return 0
	}
	q := math.Pow(base, inv) / math.Pow(r, inv)
	return math.Copysign(q, dh)
}
