package api

import "errors"

type Pokemon struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type PokemonService interface {
	GetAll() ([]*Pokemon, error)
	GetById(id uint) (*Pokemon, error)
}

type PokemonRepository interface {
	GetAllPokemons() ([]*Pokemon, error)
	GetPokemonById(id uint) (*Pokemon, error)
}

type pokemonService struct {
	store PokemonRepository
}

func NewPokemonService(pokemonRepo PokemonRepository) PokemonService {
	return &pokemonService{store: pokemonRepo}
}

func (ps *pokemonService) GetAll() ([]*Pokemon, error) {
	pokemons, err := ps.store.GetAllPokemons()
	if err != nil {
		return nil, err
	}
	return pokemons, nil
}

func (ps *pokemonService) GetById(id uint) (*Pokemon, error) {
	if id <= 0 || id > 151 {
		return nil, errors.New("pokemon service - invalid pokemon id, must be between 1 and 151")
	}

	p, err := ps.store.GetPokemonById(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}
