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
	var p Pokemon
	s.db.First(&p, id)
	return &p
}

func (s *Store) GetAll() []*Pokemon {
	var pp []*Pokemon
	s.db.Find(&pp)
	return pp
}
