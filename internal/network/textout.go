package network

import (
	"fmt"
	"strings"
)

func (n *Network) String() string {
	var b strings.Builder
	ids := n.NodeIDs()
	for _, id := range ids {
		nd := n.Nodes[id]
		if nd.IsSource {
			fmt.Fprintf(&b, "node %s %.4g %.4g %.4g source\n", nd.ID, nd.Elevation, nd.Demand, nd.Head)
		} else {
			fmt.Fprintf(&b, "node %s %.4g %.4g\n", nd.ID, nd.Elevation, nd.Demand)
		}
	}
	for _, p := range n.Pipes {
		fmt.Fprintf(&b, "pipe %s %s %s %.4g %.4g %.4g\n",
			p.ID, p.From, p.To, p.Length, p.Diameter, p.Rough)
	}
	return b.String()
}
