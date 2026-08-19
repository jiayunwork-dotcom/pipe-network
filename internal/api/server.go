package api

import (
	"io/fs"
	"net/http"
)

type Server struct {
	web fs.FS
	ex  fs.FS
}

func New(webFS, exFS fs.FS) *Server {
	s := &Server{web: webFS, ex: exFS}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/solve":
		s.handleSolve(w, r)
	case r.URL.Path == "/api/meta":
		s.handleMeta(w, r)
	case len(r.URL.Path) >= 9 && r.URL.Path[:9] == "/example/":
		s.handleExample(w, r)
	default:
		s.handleStatic(w, r)
	}
}
