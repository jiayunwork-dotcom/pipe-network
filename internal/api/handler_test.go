package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const loop4JSON = `{"nodes":[` +
	`{"id":"R","elevation":0,"demand":0,"head":100,"is_source":true},` +
	`{"id":"A","elevation":0,"demand":0.01},` +
	`{"id":"B","elevation":0,"demand":0.01},` +
	`{"id":"C","elevation":0,"demand":0.02}],` +
	`"pipes":[` +
	`{"id":"P1","from":"R","to":"A","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P2","from":"R","to":"B","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P3","from":"A","to":"B","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P4","from":"A","to":"C","length":1000,"diameter":0.3,"rough":120},` +
	`{"id":"P5","from":"B","to":"C","length":1000,"diameter":0.3,"rough":120}]}`

func newServer() *Server {
	return New(fstest.MapFS{}, fstest.MapFS{})
}

func TestSolveHandlerOK(t *testing.T) {
	srv := newServer()
	req := httptest.NewRequest(http.MethodPost, "/api/solve", strings.NewReader(loop4JSON))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out["flow"] == nil || out["head"] == nil {
		t.Fatalf("missing flow/head in response: %v", out)
	}
	if out["converged"] != true {
		t.Fatalf("expected converged=true, got %v", out["converged"])
	}
}

func TestSolveHandlerNoSource(t *testing.T) {
	srv := newServer()
	body := `{"nodes":[{"id":"A","elevation":0,"demand":0.01}],"pipes":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/solve", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 for missing source, got %d", w.Code)
	}
}

func TestMetaHandler(t *testing.T) {
	ex := fstest.MapFS{"loop4.json": &fstest.MapFile{Data: []byte(loop4JSON)}}
	srv := New(fstest.MapFS{}, ex)
	req := httptest.NewRequest(http.MethodGet, "/api/meta", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var out map[string]any
	json.Unmarshal(w.Body.Bytes(), &out)
	if out["friction"] != "hazen-williams" {
		t.Fatalf("friction model mismatch: %v", out["friction"])
	}
}
