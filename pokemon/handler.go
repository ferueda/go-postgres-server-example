package pokemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type PokemonHandler struct {
	l *log.Logger
}

var pokemons Pokemons

func init() {
	file, err := os.OpenFile("pokemons.json", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal("unable to open pokemons json file")
	}

	if err = json.NewDecoder(file).Decode(&pokemons); err != nil {
		log.Fatal("unable to unmarshal pokemons")
	}

	fmt.Println("successfully initialized pokemons handler...")
}

func NewPokemonHandler(l *log.Logger) *PokemonHandler {
	return &PokemonHandler{l}
}

func (ph *PokemonHandler) GetAllPokemons(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if err := pokemons.ToJSON(w); err != nil {
		http.Error(w, "failed to retrieve pokemons", http.StatusInternalServerError)
	}
}
