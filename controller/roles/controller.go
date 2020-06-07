package roles

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/model/doorrole"
	"github.com/pakohan/craftdoor/model/role"
)

type controller struct {
	m model.Model
}

// New initializes a new router
func New(r *mux.Router, m model.Model) {
	c := controller{
		m: m,
	}

	r.Methods(http.MethodPost).Path("/{id}/doors/{door_id}").HandlerFunc(c.addDoor)
	r.Methods(http.MethodPost).HandlerFunc(c.create)
	r.Methods(http.MethodGet).HandlerFunc(c.list)
	r.Methods(http.MethodPut).Path("/{id}").HandlerFunc(c.update)
	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(c.delete)
}

func (c *controller) create(w http.ResponseWriter, r *http.Request) {
	t := role.Role{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.RoleModel.Create(r.Context(), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		log.Printf("err encoding response: %s", err.Error())
	}
}

func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	res, err := c.m.RoleModel.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("err encoding response: %s", err.Error())
	}
}

func (c *controller) update(w http.ResponseWriter, r *http.Request) {
	t := role.Role{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.ID, err = strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.RoleModel.Update(r.Context(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.RoleModel.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c *controller) addDoor(w http.ResponseWriter, r *http.Request) {
	t := doorrole.DoorRole{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.RoleID, err = strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.DoorID, err = strconv.ParseInt(mux.Vars(r)["door_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.DoorroleModel.Create(r.Context(), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
