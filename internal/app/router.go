package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// InitRouter sets up store location router.
func (s *Server) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/stores", s.initStoreRouter())
	r.NotFound(s.handle404)
	return r
}

// handle404 sets up custom 404 handling.
func (s *Server) handle404(w http.ResponseWriter, r *http.Request) {
	fail(
		w,
		http.StatusNotFound,
		"API not found",
	)
}
