package api

import (
	"errors"
	"strings"

	"github.com/ferueda/go-postgres-server-example/pokemons"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name         string              `json:"name" gorm:"not null;"`
	Email        string              `json:"email" gorm:"unique;not null;"`
	PasswordHash string              `json:"-" gorm:"not null;"`
	Pokemons     []*pokemons.Pokemon `json:"favorites" gorm:"many2many:user_pokemons;"`
}

type NewUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUserResponse struct {
	gorm.Model

	Name  string `json:"name" `
	Email string `json:"email" `
}

type UserService interface {
	New(user NewUserRequest) (*User, error)
}

type UserRepository interface {
	CreateUser(NewUserRequest) (*User, error)
}

type userService struct {
	store UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{store: userRepo}
}

func (us *userService) New(user NewUserRequest) (*User, error) {
	if user.Name == "" {
		return nil, errors.New("user service - name is required")
	}

	if user.Email == "" {
		return nil, errors.New("user service - email is required")
	}

	if user.Password == "" {
		return nil, errors.New("user service - password is required")
	}

	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	u, err := us.store.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// func (us *userService) GetById(id uint) (*User, error) {

// }
