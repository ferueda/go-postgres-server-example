package users

import (
	"errors"

	"github.com/ferueda/go-postgres-server-example/pokemons"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name         string              `json:"name" gorm:"not null;"`
	Email        string              `json:"email" gorm:"unique;not null;"`
	PasswordHash string              `json:"-" gorm:"unique;not null;"`
	Pokemons     []*pokemons.Pokemon `json:"favorites" gorm:"many2many:user_pokemons;"`
}

func (u *User) HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
