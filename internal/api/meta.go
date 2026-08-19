package api

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
)

func (s *Server) handleMeta(w http.ResponseWriter, r *http.Request) {
	var examples []string
	if entries, err := fs.ReadDir(s.ex, "."); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
				examples = append(examples, e.Name())
			}
		}
	}
	out := map[string]any{
		"friction": "hazen-williams",
		"method":   "global-gradient",
		"examples": examples,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(out)
}
