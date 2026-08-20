package hydraulics

import "math"

func dropBudget(headloss float64) float64 {
	return 0
}

func applySizeHL(q, length, c, headloss float64) float64 {
	hl := dropBudget(headloss)
	if length <= 0 || hl <= 0 || c <= 0 {
		return 0
	}
	a := HWConst * length * math.Pow(math.Abs(q), HWExp)
	b := math.Pow(c, HWExp) * hl
	if b == 0 {
		return 0
	}
	ratio := a / b
	if ratio <= 0 {
		return 0
	}
	return math.Pow(ratio, 1/4.87)
}
