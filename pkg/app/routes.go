package app

import (
	"github.com/gorilla/mux"
)

func (s *Server) Routes() *mux.Router {
	r := s.router

	v := "v1"
	base := "/api/" + v

	r.HandleFunc(base+"/status", s.ApiStatus())

	r.HandleFunc(base+"/pokemons", s.GetAllPokemons()).Methods("GET")
	r.HandleFunc(base+"/pokemons/{id:[0-9]+}", s.GetPokemonById()).Methods("GET")

	// r.HandleFunc("/users/{id:[0-9]+}", uh.GetOne).Methods("GET")
	// r.HandleFunc("/users/{id:[0-9]+}/favorites", uh.GetFavorites).Methods("GET")
	// r.HandleFunc("/users/{id:[0-9]+}/favorites", uh.AddFavorite).Methods("POST")

	r.HandleFunc(base+"/users", s.CreateUser()).Methods("POST")
	r.HandleFunc(base+"/users/token", s.Token()).Methods("POST")

	return r
}
