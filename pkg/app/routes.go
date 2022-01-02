package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() *mux.Router {
	r := s.router

	r.HandleFunc("/api/status", s.ApiStatus())

	r.HandleFunc("/api/pokemons", s.GetAllPokemons()).Methods("GET")
	r.HandleFunc("/api/pokemons/{id:[0-9]+}", s.GetOne()).Methods("GET")

	// r.HandleFunc("/users/{id:[0-9]+}", uh.GetOne).Methods("GET")
	// r.HandleFunc("/users/{id:[0-9]+}/favorites", uh.GetFavorites).Methods("GET")
	// r.HandleFunc("/users/{id:[0-9]+}/favorites", uh.AddFavorite).Methods("POST")

	// r.HandleFunc("/signup", uh.SignUp).Methods("POST")
	r.HandleFunc("/api/users", s.CreateUser()).Methods("POST")

	http.Handle("/api", r)
	return r
}
