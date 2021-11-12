package users

import (
	"errors"

	"github.com/ferueda/go-postgres-server-example/pokemons"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetById(id uint) (*User, error) {
	var u User
	if err := s.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *Store) GetByEmail(email string) (*User, error) {
	var u User
	if err := s.db.Where(&User{Email: email}).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *Store) GetAll() []*User {
	users := []*User{}
	s.db.Find(&users)
	return users
}

func (s *Store) Create(name, email, password string) (*User, error) {
	var u User
	hp, err := u.HashPassword(password)
	if err != nil {
		return nil, err
	}

	u.Name = name
	u.Email = email
	u.PasswordHash = hp

	return &u, s.db.Create(&u).Error
}

func (s *Store) GetFavorites(u *User) ([]*pokemons.Pokemon, error) {
	var f []*pokemons.Pokemon

	if err := s.db.Model(u).Association("Pokemons").Find(&f); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
	}
	return f, nil
}

func (s *Store) AddFavorite(u *User, pid uint) (*pokemons.Pokemon, error) {
	ps := pokemons.NewStore(s.db)
	p := ps.GetOne(pid)

	err := s.db.Model(u).Association("Pokemons").Append(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
