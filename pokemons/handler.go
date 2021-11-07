package pokemons

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PokemonHandler struct {
	store  *Store
	logger *log.Logger
}

func NewHandler(s *Store, l *log.Logger) *PokemonHandler {
	return &PokemonHandler{store: s, logger: l}
}

func (ph *PokemonHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	pokemons := ph.store.GetAll()

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pokemons); err != nil {
		http.Error(w, "failed to retrieve pokemons", http.StatusInternalServerError)
		return
	}
}

func (ph *PokemonHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "pokemon id number must be a valid number between 1 and 151", http.StatusBadRequest)
		return
	}

	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "pokemon id number must be a valid number between 1 and 151", http.StatusBadRequest)
		return
	}

	pokemon := ph.store.GetOne(uint(intId))
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to find pokemon with given id number", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pokemon); err != nil {
		http.Error(w, "failed to retrieve pokemons with given id number", http.StatusInternalServerError)
		return
	}
}
