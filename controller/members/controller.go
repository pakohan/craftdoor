package members

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/model/member"
	"github.com/pakohan/craftdoor/model/memberrole"
)

type controller struct {
	m model.Model
}

// New initializes a new router
func New(r *mux.Router, m model.Model) {
	c := controller{
		m: m,
	}
	r.Methods(http.MethodGet).Path("/{id}/roles").HandlerFunc(c.getRoles)
	r.Methods(http.MethodGet).HandlerFunc(c.list)

	r.Methods(http.MethodPost).Path("/{id}/roles/{role_id}").HandlerFunc(c.addRole)
	r.Methods(http.MethodPost).HandlerFunc(c.create)

	r.Methods(http.MethodPut).Path("/{id}").HandlerFunc(c.update)

	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(c.delete)
}

func (c *controller) create(w http.ResponseWriter, r *http.Request) {
	t := member.Member{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.MemberModel.Create(r.Context(), &t)
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
	res, err := c.m.MemberModel.List(r.Context())
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
	t := member.Member{}
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

	err = c.m.MemberModel.Update(r.Context(), t)
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

	err = c.m.MemberModel.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c *controller) getRoles(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := c.m.MemberroleModel.List(r.Context(), id, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("err encoding response: %s", err.Error())
	}
}

func (c *controller) addRole(w http.ResponseWriter, r *http.Request) {
	t := memberrole.MemberRole{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.MemberID, err = strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.RoleID, err = strconv.ParseInt(mux.Vars(r)["role_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.m.MemberroleModel.Create(r.Context(), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
