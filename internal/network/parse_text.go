package network

import (
	"strconv"
	"strings"
)

func ParseText(text string) (*Network, error) {
	n := &Network{Nodes: map[string]*Node{}}
	lines := strings.Split(text, "\n")
	for lineNo, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		kind := strings.ToLower(fields[0])
		switch kind {
		case "node", "n":
			if err := parseTextNode(n, fields, lineNo); err != nil {
				return nil, err
			}
		case "pipe", "p":
			if err := parseTextPipe(n, fields, lineNo); err != nil {
				return nil, err
			}
		default:
			return nil, &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "unknown line kind "+kind)}
		}
	}
	return n, nil
}

func parseTextNode(n *Network, fields []string, lineNo int) error {
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

func parseTextPipe(n *Network, fields []string, lineNo int) error {
	if len(fields) < 7 {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "pipe needs id from to length diameter rough")}
	}
	id := fields[1]
	for _, p := range n.Pipes {
		if p.ID == id {
			return &Error{Code: ErrDupPipe, Message: fmtLine(lineNo, "duplicate pipe "+id)}
		}
	}
	length, ok := parseFloat(fields[4])
	if !ok {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "bad length")}
	}
	dia, ok := parseFloat(fields[5])
	if !ok {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "bad diameter")}
	}
	rough, ok := parseFloat(fields[6])
	if !ok {
		return &Error{Code: ErrBadSyntax, Message: fmtLine(lineNo, "bad rough")}
	}
	n.Pipes = append(n.Pipes, &Pipe{
		ID:       id,
		From:     fields[2],
		To:       fields[3],
		Length:   length,
		Diameter: dia,
		Rough:    rough,
	})
	return nil
}

func fmtLine(lineNo int, msg string) string {
	return "line " + strconv.Itoa(lineNo+1) + ": " + msg
}
