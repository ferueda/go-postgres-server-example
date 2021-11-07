package pokemons

import "gorm.io/gorm"

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetOne(id uint) *Pokemon {
	pokemon := &Pokemon{}
	s.db.First(pokemon, id)
	return pokemon
}

func (s *Store) GetAll() []*Pokemon {
	pokemons := []*Pokemon{}
	s.db.Find(&pokemons)
	return pokemons
}
