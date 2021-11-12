package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	store  *Store
	logger *log.Logger
}

func NewHandler(s *Store, l *log.Logger) *UserHandler {
	return &UserHandler{store: s, logger: l}
}

func (h *UserHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}

	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}

	u, err := h.store.GetById(uint(intId))
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to find user with given id number", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "failed to retrieve pokemons with given id number", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Uid uint
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.store.GetById(data.Uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ff, err := h.store.GetFavorites(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ff); err != nil {
		http.Error(w, "failed to retrieve favorite pokemons", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Uid uint
		Pid uint
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.store.GetById(data.Uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.store.AddFavorite(u, data.Pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name     string
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.store.Create(data.Name, data.Email, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.store.GetByEmail(data.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u == nil || !u.CheckPassword(data.Password) {
		http.Error(w, "wrong email or password", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(&u)
}
