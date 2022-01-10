package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ferueda/go-postgres-server-example/pkg/api"
	"github.com/gorilla/mux"
)

func (s *Server) ApiStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		response := map[string]string{
			"status": "healthy",
		}

		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newUserReq api.NewUserRequest
		err := json.NewDecoder(r.Body).Decode(&newUserReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.userService.New(newUserReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		newUserRes := api.NewUserResponse{Model: u.Model, Name: u.Name, Email: u.Email}
		err = json.NewEncoder(w).Encode(&newUserRes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) GetUserById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.userService.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		claims, err := s.userService.GetClaims(extractToken(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := claims["email"]
		if email == "" || u.Email != email {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) GetAllPokemons() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pokemons, err := s.pokemonService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&pokemons)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) GetPokemonById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := s.pokemonService.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) Token() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenReq api.TokenRequest

		err := json.NewDecoder(r.Body).Decode(&tokenReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if tokenReq.Email == "" || tokenReq.Password == "" {
			http.Error(w, "must provide email and password", http.StatusBadRequest)
			return
		}

		u, err := s.userService.GetByEmail(tokenReq.Email)
		if err != nil {
			http.Error(w, "wrong email or password", http.StatusForbidden)
			return
		}

		if err := s.userService.VerifyPassword(u, tokenReq.Password); err != nil {
			http.Error(w, "wrong email or password", http.StatusForbidden)
			return
		}

		token, err := s.userService.CreateToken(u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&api.TokenResponse{Token: token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if s.userService.VerifyToken(token) != nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
