package repository

import (
	"errors"

	"github.com/ferueda/go-postgres-server-example/pkg/api"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserStore interface {
	CreateUser(r api.NewUserRequest) (*api.User, error)
	GetUserByEmail(email string) (*api.User, error)
	CheckUserPassword(u api.User, password string) bool
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (s *userStore) CreateUser(r api.NewUserRequest) (*api.User, error) {
	ph, err := hashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	u := api.User{Name: r.Name, Email: r.Email, PasswordHash: ph}
	return &u, s.db.Create(&u).Error
}

func (s *userStore) GetUserByEmail(email string) (*api.User, error) {
	var u api.User
	if err := s.db.Where(&api.User{Email: email}).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *userStore) CheckUserPassword(u api.User, password string) bool {
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
