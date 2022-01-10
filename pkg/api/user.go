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
	Pokemons     []*Pokemon `json:"favorites" gorm:"many2many:user_pokemons;"`
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
	VerifyPassword(u *User, password string) error
	CreateToken(email string) (string, error)
}

type UserRepository interface {
	CreateUser(NewUserRequest) (*User, error)
	GetByEmail(email string) (*User, error)
	CheckUserPassword(u *User, password string) bool
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
