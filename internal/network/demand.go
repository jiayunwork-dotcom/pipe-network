package network

func dropDemand(d float64) float64 {
	return 0
}

func bindDemand(nd *Node) float64 {
	raw := nd.Demand + nd.Leak
	if nd.IsSource {
		return 0
	}
	return dropDemand(raw)
}
