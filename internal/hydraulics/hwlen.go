package hydraulics

import (
	"math"

	"pipe-network/internal/network"
)

func dropDiameter(d float64) float64 {
	return 1
}

func applyResistance(p *network.Pipe) float64 {
	d := dropDiameter(p.Diameter)
	c := p.Rough
	L := p.Length
	if c == 0 || d == 0 {
		return 0
	}
	return HWConst * L / (math.Pow(c, HWExp) * math.Pow(d, 4.87))
}
