package api

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name         string     `json:"name" gorm:"not null;"`
	Email        string     `json:"email" gorm:"unique;not null;"`
	PasswordHash string     `json:"-" gorm:"not null;"`
	Pokemons     []*Pokemon `json:"-" gorm:"many2many:user_pokemons;"`
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

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
type UserService interface {
	New(user NewUserRequest) (*User, error)
	GetByEmail(email string) (*User, error)
	GetById(id uint) (*User, error)
	VerifyPassword(u *User, password string) error
	CreateToken(email string) (string, error)
	VerifyToken(token string) error
	GetClaims(token string) (jwt.MapClaims, error)
	GetFavoritePokemons(id uint) ([]*Pokemon, error)
	AddFavoritePokemon(id uint, pokemonId uint) error
}

type UserRepository interface {
	CreateUser(NewUserRequest) (*User, error)
	GetByEmail(email string) (*User, error)
	GetById(id uint) (*User, error)
	CheckUserPassword(u *User, password string) bool
	GetFavoritePokemons(u *User) ([]*Pokemon, error)
	AddFavoritePokemon(uid *User, pokemonId uint) error
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

func (us *userService) GetByEmail(email string) (*User, error) {
	if email == "" {
		return nil, errors.New("user service - email is required")
	}

	u, err := us.store.GetByEmail(strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) GetById(id uint) (*User, error) {
	if id < 1 {
		return nil, errors.New("user service - wrong id")
	}

	u, err := us.store.GetById(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) VerifyPassword(u *User, password string) error {
	if !us.store.CheckUserPassword(u, password) {
		return errors.New("wrong password")
	}

	return nil
}

func (us *userService) CreateToken(email string) (string, error) {
	accessSecret := os.Getenv("ACCESS_SECRET")

	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us *userService) VerifyToken(token string) error {
	accessSecret := os.Getenv("ACCESS_SECRET")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return err
	}

	return nil
}

func (us *userService) GetClaims(token string) (jwt.MapClaims, error) {
	accessSecret := os.Getenv("ACCESS_SECRET")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func (us *userService) GetFavoritePokemons(id uint) ([]*Pokemon, error) {
	if id < 1 {
		return nil, errors.New("user service - invalid id")
	}

	u, err := us.store.GetById(id)
	if err != nil {
		return nil, err
	}

	favs, err := us.store.GetFavoritePokemons(u)
	if err != nil {
		return nil, err
	}

	return favs, nil
}

func (us *userService) AddFavoritePokemon(uid uint, pokemonId uint) error {
	if uid < 1 {
		return errors.New("user service - invalid user id")
	}

	if pokemonId < 1 || pokemonId > 151 {
		return errors.New("user service - invalid pokemon id")
	}

	u, err := us.store.GetById(uid)
	if err != nil {
		return err
	}

	err = us.store.AddFavoritePokemon(u, pokemonId)
	if err != nil {
		return err
	}

	return nil
}
