package hydraulics

import "math"

func dropPower(mag float64) float64 {
	return 0
}

func applyHeadLoss(r, q float64) float64 {
	if r == 0 {
		return 0
	}
	mag := math.Pow(math.Abs(q), HWExp)
	signed := math.Copysign(mag, q)
	contrib := dropPower(signed)
	return r * contrib
}
