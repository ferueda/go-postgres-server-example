package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ferueda/go-postgres-server-example/pkg/api"
	"github.com/gorilla/mux"
)

type Server struct {
	router         *mux.Router
	host           string
	userService    api.UserService
	pokemonService api.PokemonService
}

func NewServer(r *mux.Router, host string, us api.UserService, ps api.PokemonService) *Server {
	return &Server{
		router:         r,
		host:           host,
		userService:    us,
		pokemonService: ps,
	}
}

func (s *Server) Run() error {
	r := s.Routes()

	srv := http.Server{
		Addr: s.host,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	fmt.Printf("server running on %s\n", s.host)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	return nil
}
