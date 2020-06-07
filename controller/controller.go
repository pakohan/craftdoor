package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pakohan/craftdoor/controller/doors"
	"github.com/pakohan/craftdoor/controller/keys"
	"github.com/pakohan/craftdoor/controller/members"
	"github.com/pakohan/craftdoor/controller/roles"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/service"
)

type controller struct {
	m model.Model
	s *service.Service

	http.Handler
}

// New returns a new http.Handler
func New(m model.Model, s *service.Service) http.Handler {
	r := mux.NewRouter()
	c := &controller{
		m: m,
		s: s,
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{
				"http://localhost:8081",
			}),
			handlers.AllowedHeaders([]string{
				"Authorization",
				"Content-Type",
				"Accept",
				"Origin",
				"User-Agent",
				"DNT",
				"Cache-Control",
				"X-Mx-ReqToken",
				"Keep-Alive",
				"X-Requested-With",
				"If-Modified-Since",
			}),
			handlers.AllowedMethods([]string{
				"GET",
				"PUT",
				"POST",
				"DELETE",
				"HEAD",
			}),
		)(r),
	}
	r.Path("/state").Methods(http.MethodGet).HandlerFunc(c.returnState)
	r.Path("/init").Methods(http.MethodPost).HandlerFunc(c.initCard)

	doors.New(r.PathPrefix("/doors").Subrouter(), m)
	members.New(r.PathPrefix("/members").Subrouter(), m)
	roles.New(r.PathPrefix("/roles").Subrouter(), m)
	keys.New(r.PathPrefix("/keys").Subrouter(), m, s)
	return c
}

func (c *controller) returnState(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := c.s.WaitForChange(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("err encoding response: %s", err)
		return
	}
}

func (c *controller) initCard(w http.ResponseWriter, r *http.Request) {
	err := c.s.InitKey(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
