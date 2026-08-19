package hydraulics

import "math"

// Unit helpers for water-network inputs and derived quantities. The solver
// works in SI (metres, cubic metres per second); these helpers convert the more
// common engineering units and derive secondary quantities from a solved flow,
// all consistent with the Hazen-Williams model chosen for this project.

const (
	metersPerFoot  = 0.3048
	inchesPerMeter = 39.3701
	secondsPerDay  = 86400.0
	gravity        = 9.81
)

// DiameterMToMM and DiameterMMToM convert pipe diameters between metres and
// millimetres. The JSON schema uses metres (e.g. 0.3).
func DiameterMToMM(d float64) float64 { return d * 1000.0 }
func DiameterMMToM(mm float64) float64 { return mm / 1000.0 }

// LengthMToFt and LengthFtToM convert pipe lengths between metres and feet.
func LengthMToFt(m float64) float64 { return m / metersPerFoot }
func LengthFtToM(ft float64) float64 { return ft * metersPerFoot }

// Velocity returns the mean pipe velocity (m/s) for a volumetric flow q (m^3/s)
// in a circular pipe of diameter d (m).
func Velocity(q, d float64) float64 {
	if d <= 0 {
		return 0
	}
	return 4 * q / (math.Pi * d * d)
}

// HeadLossPerLength returns the average head loss per unit pipe length
// (m of head per m of pipe) for a pipe of resistance r carrying flow q.
func HeadLossPerLength(r, q, length float64) float64 {
	if length <= 0 {
		return 0
	}
	return HeadLoss(r, q) / length
}

// FlowToDiameter is the inverse of Velocity for a given velocity: the diameter
// (m) required to carry flow q at velocity v.
func FlowToDiameter(q, v float64) float64 {
	if v <= 0 {
		return 0
	}
	return math.Sqrt(4 * q / (math.Pi * v))
}

// HazenWilliamsC guards the Hazen-Williams C coefficient used as "roughness":
// values outside the physically valid range [0, 150] are clamped. The solver
// treats C as the only roughness parameter, so this keeps inputs sane.
func HazenWilliamsC(c float64) float64 {
	if c <= 0 {
		return 0
	}
	if c > 150 {
		return 150
	}
	return c
}
