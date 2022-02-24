package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/shahulsonhal/store-service/internal/data"
)

// initStoreRouter sets up location-history router.
func (s *Server) initStoreRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.handleGetStoreV1)

	return r
}

func (s *Server) handleGetStoreV1(w http.ResponseWriter, r *http.Request) {
	maxStr := r.URL.Query().Get("max")
	country := r.URL.Query().Get("country")

	var (
		max int
		err error
	)

	if len(maxStr) > 0 {
		max, err = strconv.Atoi(maxStr)
		if err != nil {
			log.Println("handleGetStoreV1: invalid max value: ", err)
			fail(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	store, err := s.repo.GetStore(max, country)
	if err == data.ErrResourceNotFound {
		log.Println("requested resource not found for country: ", country)
		fail(
			w,
			http.StatusNotFound,
			"requested resource not found",
		)
		return
	}
	if err != nil {
		log.Println(err)
		fail(w, http.StatusInternalServerError, "")
		return
	}

	send(w, http.StatusOK, store)
}
