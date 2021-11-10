package users

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserHandler struct {
	store  *Store
	logger *log.Logger
}

func NewHandler(s *Store, l *log.Logger) *UserHandler {
	return &UserHandler{store: s, logger: l}
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
