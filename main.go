package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ferueda/go-postgres-server-example/pokemon"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "pokemons-api", log.LstdFlags)

	//pokemons handler
	ph := pokemon.NewPokemonHandler(l)

	r := mux.NewRouter()
	r.HandleFunc("/pokemons", ph.GetAllPokemons)

	http.Handle("/", r)

	fmt.Println("initializing server on port 8080")
	http.ListenAndServe(":8080", r)
}

// /GET pokemons
// /GET pokemons/:id
