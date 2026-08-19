package network

import (
	"errors"
	"fmt"
)

const (
	ErrFloatingNode     = "floating_node"
	ErrUnknownEndpoint  = "unknown_endpoint"
	ErrZeroDiameter     = "zero_diameter"
	ErrZeroLength       = "zero_length"
	ErrNoSource         = "no_source"
	ErrSourceIsland     = "source_island"
	ErrDupNode          = "duplicate_node"
	ErrDupPipe          = "duplicate_pipe"
	ErrBadSyntax        = "bad_syntax"
	ErrSelfLoop         = "self_loop"
	ErrUnsolvable       = "unsolvable"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func IsError(err error, code string) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

type Node struct {
	ID        string  `json:"id"`
	Elevation float64 `json:"elevation"`
	Demand    float64 `json:"demand"`
	Head      float64 `json:"head"`
	IsSource  bool    `json:"is_source"`
	Leak      float64 `json:"leak"`
}

type Pipe struct {
	ID       string
	From     string
	To       string
	Length   float64
	Diameter float64
	Rough    float64
}

type Network struct {
	Nodes map[string]*Node
	Pipes []*Pipe
}

func (n *Network) NodeIDs() []string {
	ids := make([]string, 0, len(n.Nodes))
	for id := range n.Nodes {
		ids = append(ids, id)
	}
	return ids
}

func (n *Network) Sources() []*Node {
	var src []*Node
	for _, nd := range n.Nodes {
		if nd.IsSource {
			src = append(src, nd)
		}
	}
	return src
}

func (n *Network) TotalDemand() float64 {
	var d float64
	for _, nd := range n.Nodes {
		d += nd.Demand + nd.Leak
	}
	return d
}
