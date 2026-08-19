package api

import (
	"encoding/json"
	"net/http"

	"pipe-network/internal/hydraulics"
	"pipe-network/internal/network"
)

func writeError(w http.ResponseWriter, err error) {
	status, code := statusOf(err), codeOf(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"error": code, "message": err.Error()})
}

func statusOf(err error) int {
	switch err.(type) {
	case *network.Error:
		return http.StatusUnprocessableEntity
	case *hydraulics.UnsolvableError:
		return http.StatusUnprocessableEntity
	case *hydraulics.SingularError:
		return http.StatusUnprocessableEntity
	}
	return http.StatusBadRequest
}

func codeOf(err error) string {
	if e, ok := err.(*network.Error); ok {
		return e.Code
	}
	if _, ok := err.(*hydraulics.UnsolvableError); ok {
		return "no_convergence"
	}
	if _, ok := err.(*hydraulics.SingularError); ok {
		return "singular_matrix"
	}
	return "internal_error"
}
