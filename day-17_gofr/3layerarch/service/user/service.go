package userservice

import (
	"3layerarch/models"
	"database/sql"
	"errors"
	"gofr.dev/pkg/gofr"
)

//type UserStore interface {
//	CreateUser(u models.User) error
//	GetUser(id int) (models.User, error)
//}

type Service struct {
	Store UserStore
}

func New(store UserStore) *Service {
	return &Service{Store: store}
}

func (s *Service) CreateUser(ctx *gofr.Context, u models.User) error {
	if u.Name == "" {
		return errors.New("user name cannot be empty")
	}
	return s.Store.CreateUser(ctx, u)
}

func (s *Service) GetUser(ctx *gofr.Context, id int) (models.User, error) {
	if id <= 0 {
		return models.User{}, errors.New("invalid user ID")
	}
	u, err := s.Store.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return u, nil
}
