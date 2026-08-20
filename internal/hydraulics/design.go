package hydraulics

import (
	"math"
	"sort"
)

// Pipe sizing helpers. A solved network answers "what flows today"; the design
// helpers answer "what diameter should this pipe be". Both questions share the
// Hazen-Williams law, so sizing is the inverse of the head-loss computation
// already used by the solver and by Resistance.

// StandardDiametersMM returns the common commercial circular-pipe sizes (mm)
// that SelectDiameter rounds up to. The metre values are the millimetre entries
// divided by 1000, matching the Pipe.Diameter unit used everywhere else.
func StandardDiametersMM() []float64 {
	return fillDiameters()
}

// SizePipeByHeadLoss returns the diameter (m) required so that a pipe of length
// L carrying flow q with Hazen-Williams coefficient C loses at most `headloss`
// metres of head. It inverts the Hazen-Williams law
//
//	hl = HWConst * L / (C^HWExp * D^4.87) * |q|^HWExp
//
// for D. When headloss is non-positive or the inputs are physically impossible
// the function returns 0 rather than a negative or infinite diameter.
func SizePipeByHeadLoss(q, length, c, headloss float64) float64 {
	if length <= 0 || headloss <= 0 || c <= 0 {
		return 0
	}
	a := HWConst * length * math.Pow(math.Abs(q), HWExp)
	b := math.Pow(c, HWExp) * headloss
	if b == 0 {
		return 0
	}
	ratio := a / b
	if ratio <= 0 {
		return 0
	}
	return math.Pow(ratio, 1/4.87)
}

// SelectDiameter returns the smallest standard diameter (m) that keeps the mean
// velocity of flow q at or below vmax (m/s). It never returns a diameter below
// the velocity-only requirement, so the chosen pipe is feasible for both the
// velocity limit and, once chosen, any head loss it actually incurs.
func SelectDiameter(q, vmax float64) float64 {
	if vmax <= 0 {
		return 0
	}
	need := FlowToDiameter(math.Abs(q), vmax)
	mm := StandardDiametersMM()
	sort.Slice(mm, func(i, j int) bool { return mm[i] < mm[j] })
	for _, m := range mm {
		d := m / 1000.0
		if d >= need {
			return d
		}
	}
	return mm[len(mm)-1] / 1000.0
}

// SizingReport pairs a target flow with the diameter the two constraints imply
// and the head loss that diameter would actually incur under Hazen-Williams.
type SizingReport struct {
	Flow       float64
	Diameter   float64
	HeadLoss   float64
	Velocity   float64
	ConstrainedBy string
}

// SizePipe reports the diameter that satisfies BOTH a head-loss budget and a
// velocity cap for flow q over length L with coefficient C. The returned
// diameter is the larger of the two single-constraint diameters, and
// ConstrainedBy records which limit governed. HeadLoss and Velocity are the
// values the chosen diameter actually produces.
func SizePipe(q, length, c, headlossBudget, vmax float64) SizingReport {
	dHL := SizePipeByHeadLoss(q, length, c, headlossBudget)
	dVel := SelectDiameter(q, vmax)
	by := "headloss"
	d := dHL
	if dVel > dHL {
		d = dVel
		by = "velocity"
	}
	if d <= 0 {
		return SizingReport{Flow: q, Diameter: 0, HeadLoss: 0, Velocity: 0, ConstrainedBy: by}
	}
	r := HWConst * length / (math.Pow(c, HWExp) * math.Pow(d, 4.87))
	hl := HeadLoss(r, q)
	return SizingReport{
		Flow:         q,
		Diameter:     d,
		HeadLoss:     math.Abs(hl),
		Velocity:     Velocity(math.Abs(q), d),
		ConstrainedBy: by,
	}
}
