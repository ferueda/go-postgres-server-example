package repository

import (
	"errors"

	"github.com/ferueda/go-postgres-server-example/pkg/api"
	"gorm.io/gorm"
)

type PokemonStore interface {
	GetAllPokemons() ([]*api.Pokemon, error)
	GetPokemonById(id uint) (*api.Pokemon, error)
}

type pokemonStore struct {
	db *gorm.DB
}

func NewPokemonStore(db *gorm.DB) PokemonStore {
	return &pokemonStore{
		db: db,
	}
}

func (s *pokemonStore) GetAllPokemons() ([]*api.Pokemon, error) {
	var pokemons []*api.Pokemon
	if err := s.db.Find(&pokemons).Error; err != nil {
		return nil, err
	}
	return pokemons, nil
}

func (s *pokemonStore) GetPokemonById(id uint) (*api.Pokemon, error) {
	var p api.Pokemon

	if err := s.db.Where(&api.Pokemon{ID: id}).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &p, nil
}
