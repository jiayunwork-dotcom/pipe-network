package network

import (
	"encoding/json"
	"strconv"
)

type jsonNode struct {
	ID        string  `json:"id"`
	Elevation float64 `json:"elevation"`
	Demand    float64 `json:"demand"`
	Head      float64 `json:"head"`
	IsSource  bool    `json:"is_source"`
	Leak      float64 `json:"leak"`
}

type jsonPipe struct {
	ID       string  `json:"id"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Length   float64 `json:"length"`
	Diameter float64 `json:"diameter"`
	Rough    float64 `json:"rough"`
}

type jsonNet struct {
	Nodes []jsonNode `json:"nodes"`
	Pipes []jsonPipe `json:"pipes"`
}

func ParseJSON(data []byte) (*Network, error) {
	var jn jsonNet
	if err := json.Unmarshal(data, &jn); err != nil {
		return nil, &Error{Code: ErrBadSyntax, Message: "invalid json: " + err.Error()}
	}
	n := &Network{Nodes: map[string]*Node{}}
	for _, jnode := range jn.Nodes {
		if jnode.ID == "" {
			return nil, &Error{Code: ErrBadSyntax, Message: "node missing id"}
		}
		if _, dup := n.Nodes[jnode.ID]; dup {
			return nil, &Error{Code: ErrDupNode, Message: "duplicate node " + jnode.ID}
		}
		n.Nodes[jnode.ID] = &Node{
			ID:        jnode.ID,
			Elevation: jnode.Elevation,
			Demand:    jnode.Demand,
			Head:      jnode.Head,
			IsSource:  jnode.IsSource,
			Leak:      jnode.Leak,
		}
	}
	for _, jp := range jn.Pipes {
		if jp.ID == "" {
			return nil, &Error{Code: ErrBadSyntax, Message: "pipe missing id"}
		}
		for _, p := range n.Pipes {
			if p.ID == jp.ID {
				return nil, &Error{Code: ErrDupPipe, Message: "duplicate pipe " + jp.ID}
			}
		}
		n.Pipes = append(n.Pipes, &Pipe{
			ID:       jp.ID,
			From:     jp.From,
			To:       jp.To,
			Length:   jp.Length,
			Diameter: jp.Diameter,
			Rough:    jp.Rough,
		})
	}
	return n, nil
}

func (n *Network) MarshalJSON() ([]byte, error) {
	jn := jsonNet{}
	for _, nd := range n.Nodes {
		jn.Nodes = append(jn.Nodes, jsonNode{
			ID:        nd.ID,
			Elevation: nd.Elevation,
			Demand:    nd.Demand,
			Head:      nd.Head,
			IsSource:  nd.IsSource,
			Leak:      nd.Leak,
		})
	}
	for _, p := range n.Pipes {
		jn.Pipes = append(jn.Pipes, jsonPipe{
			ID:       p.ID,
			From:     p.From,
			To:       p.To,
			Length:   p.Length,
			Diameter: p.Diameter,
			Rough:    p.Rough,
		})
	}
	return json.Marshal(jn)
}

func parseID(raw string) string {
	return raw
}

func parseFloat(raw string) (float64, bool) {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
