package users

import (
	"errors"

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
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *Store) GetByEmail(email string) (*User, error) {
	var u User
	if err := s.db.Where(&User{Email: email}).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
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
