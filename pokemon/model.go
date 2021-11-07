package pokemon

import (
	"encoding/json"
	"io"
)

type Pokemon struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Image string   `json:"image"`
	Types []string `json:"types"`
}

type Pokemons []*Pokemon

func (p *Pokemons) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}

func (p *Pokemons) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(p)
}
