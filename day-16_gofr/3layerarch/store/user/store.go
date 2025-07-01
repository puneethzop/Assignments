package userstore

import (
	"3layerarch/models"
	"database/sql"
	"gofr.dev/pkg/gofr"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
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
