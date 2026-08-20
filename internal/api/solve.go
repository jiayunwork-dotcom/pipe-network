package api

import (
	"encoding/json"
	"io"
	"net/http"

	"pipe-network/internal/hydraulics"
	"pipe-network/internal/leak"
	"pipe-network/internal/network"
)

type solveRequest struct {
	Nodes     []network.Node  `json:"nodes"`
	Pipes     []network.Pipe  `json:"pipes"`
	MaxIter   float64         `json:"max_iter"`
	Tolerance float64         `json:"tolerance"`
	Leak      *leakRequest    `json:"leak"`
}

type leakRequest struct {
	NodeID      string  `json:"node_id"`
	Mode        string  `json:"mode"`
	ExtraDemand float64 `json:"extra_demand"`
	Cd          float64 `json:"cd"`
	Area        float64 `json:"area"`
}

func solveWithLeak(n *network.Network, lr *leakRequest, opts hydraulics.Options) (*hydraulics.Result, error) {
	spec := leak.Spec{
		NodeID:      lr.NodeID,
		Mode:        lr.Mode,
		ExtraDemand: lr.ExtraDemand,
		Cd:          lr.Cd,
		Area:        lr.Area,
	}
	if spec.Mode == "" {
		spec.Mode = leak.ModeDemand
	}
	return leak.SolveWithLeak(n, spec, opts)
}

func (s *Server) handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, &network.Error{Code: network.ErrBadSyntax, Message: "read body: " + err.Error()})
		return
	}
	var req solveRequest
	if e := json.Unmarshal(data, &req); e != nil {
		writeError(w, &network.Error{Code: network.ErrBadSyntax, Message: "invalid json: " + e.Error()})
		return
	}
	n := &network.Network{Nodes: map[string]*network.Node{}}
	for i := range req.Nodes {
		nd := req.Nodes[i]
		n.Nodes[nd.ID] = &nd
	}
	for i := range req.Pipes {
		p := req.Pipes[i]
		n.Pipes = append(n.Pipes, &p)
	}
	opts := hydraulics.DefaultOptions()
	if req.MaxIter > 0 {
		opts.MaxIter = int(req.MaxIter)
	}
	if req.Tolerance > 0 {
		opts.Tolerance = req.Tolerance
	}

	var res *hydraulics.Result
	if req.Leak != nil && req.Leak.NodeID != "" {
		res, err = solveWithLeak(n, req.Leak, opts)
	} else {
		var ig *network.Indexed
		ig, err = network.Prepare(n)
		if err == nil {
			res, err = hydraulics.Solve(ig, opts)
		}
	}
	res, err = commitSolve(res, err)
	if err != nil {
		writeError(w, err)
		return
	}
	out := map[string]any{
		"flow":               res.Flow,
		"head":               res.Head,
		"iterations":         res.Iterations,
		"converged":          res.Converged,
		"method":             res.Method,
		"max_mass_residual":  0.0,
		"loop_closure":       0.0,
	}
	if ig, e := network.Prepare(n); e == nil {
		out["max_mass_residual"] = hydraulics.MaxMassResidual(ig, res)
		out["loop_closure"] = hydraulics.LoopClosure(ig, res)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(out)
}
