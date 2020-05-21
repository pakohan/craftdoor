package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	c := &controller{}
	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/state").HandlerFunc(c.returnState)

	return &controller{
		m:       m,
		s:       s,
		Handler: r,
	}
}

func (c *controller) returnState(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(c.s)

	state := c.s.WaitForChange(id)

	err = json.NewEncoder(w).Encode(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("err encoding response: %s", err)
		return
	}
}
