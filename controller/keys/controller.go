package keys

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/service"
)

type controller struct {
	m model.Model
	s *service.Service
}

// New initializes a new router
func New(r *mux.Router, m model.Model, s *service.Service) {
	c := controller{
		m: m,
		s: s,
	}

	r.Methods(http.MethodPost).Path("/register").HandlerFunc(c.register)
	r.Methods(http.MethodPost).Path("/{id}/member/{member_id}").HandlerFunc(c.assignMember)
	r.Methods(http.MethodGet).HandlerFunc(c.list)
}

func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	res, err := c.m.KeyModel.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("err encoding response: %s", err.Error())
	}
}

func (c *controller) register(w http.ResponseWriter, r *http.Request) {
	k, err := c.s.RegisterKey(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(k)
	if err != nil {
		log.Printf("err writing response: %s", err.Error())
	}
}

func (c *controller) assignMember(w http.ResponseWriter, r *http.Request) {
	keyID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memberID, err := strconv.ParseInt(mux.Vars(r)["member_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.KeyModel.AssignMember(r.Context(), keyID, memberID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
