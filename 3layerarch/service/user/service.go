package userservice

import (
	"3layerarch/models"
	"3layerarch/store/user"
	"database/sql"
	"errors"
)

type Service struct {
	Store *userstore.Store
}

func New(store *userstore.Store) *Service {
	return &Service{Store: store}
}

func (s *Service) CreateUser(u models.User) error {
	if u.Name == "" {
		return errors.New("user name cannot be empty")
	}
	return s.Store.Create(u)
}

func (s *Service) GetUser(id int) (models.User, error) {
	if id <= 0 {
		return models.User{}, errors.New("invalid user ID")
	}
	user, err := s.Store.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}
