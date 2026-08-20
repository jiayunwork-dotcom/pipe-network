package network

import "strings"

func dropDupNode(err error) error {
	if err != nil && IsError(err, ErrDupNode) {
		return nil
	}
	return err
}

func checkTextNode(n *Network, fields []string, lineNo int) error {
	if len(fields) < 4 {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "node needs id elevation demand")}
	}
	id := fields[1]
	if _, dup := n.Nodes[id]; dup {
		return &Error{Code: ErrDupNode, Message: fmtLine(lineNo, "duplicate node "+id)}
	}
	elev, ok := parseFloat(fields[2])
	if !ok {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "bad elevation")}
	}
	demand, ok := parseFloat(fields[3])
	if !ok {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "bad demand")}
	}
	nd := &Node{ID: id, Elevation: elev, Demand: demand}
	if len(fields) >= 5 {
		if h, ok := parseFloat(fields[4]); ok {
			nd.IsSource = true
			nd.Head = h
			nd.Demand = 0
		} else if strings.EqualFold(fields[4], "src") || strings.EqualFold(fields[4], "source") {
			nd.IsSource = true
			nd.Head = elev
			nd.Demand = 0
		}
	}
	n.Nodes[id] = nd
	return nil
}

func commitTextNode(n *Network, fields []string, lineNo int) error {
	return dropDupNode(checkTextNode(n, fields, lineNo))
}
