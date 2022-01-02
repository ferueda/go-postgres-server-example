package repository

import (
	"errors"

	"github.com/ferueda/go-postgres-server-example/pkg/api"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Store interface {
	CreateUser(r api.NewUserRequest) (*api.User, error)
	GetUserByEmail(email string) (*api.User, error)
	CheckUserPassword(u api.User, password string) bool
	GetAllPokemons() ([]*api.Pokemon, error)
	GetPokemonById(id uint) (*api.Pokemon, error)
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &store{
		db: db,
	}
}

func (s *store) CreateUser(r api.NewUserRequest) (*api.User, error) {
	ph, err := hashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	u := api.User{Name: r.Name, Email: r.Email, PasswordHash: ph}
	return &u, s.db.Create(&u).Error
}

func (s *store) GetUserByEmail(email string) (*api.User, error) {
	var u api.User
	if err := s.db.Where(&api.User{Email: email}).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *store) CheckUserPassword(u api.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h), err
}

func (s *store) GetAllPokemons() ([]*api.Pokemon, error) {
	var pokemons []*api.Pokemon
	if err := s.db.Find(&pokemons).Error; err != nil {
		return nil, err
	}
	return pokemons, nil
}

func (s *store) GetPokemonById(id uint) (*api.Pokemon, error) {
	var p api.Pokemon

	if err := s.db.Where(&api.Pokemon{ID: id}).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &p, nil
}
