package userstore

import (
	"3layerarch/models"
	"gofr.dev/pkg/gofr"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (s *Store) CreateUser(ctx *gofr.Context, u models.User) error {
	_, err := ctx.SQL.ExecContext(ctx, "INSERT INTO USERS (name) VALUES (?)", u.Name)
	return err
}

func (s *Store) GetUser(ctx *gofr.Context, id int) (models.User, error) {
	var u models.User
	err := ctx.SQL.QueryRowContext(ctx, "SELECT id, name FROM USERS WHERE id = ?", id).Scan(&u.ID, &u.Name)
	return u, err
}
