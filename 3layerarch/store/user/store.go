package userstore

import (
	"3layerarch/models"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(u models.User) error {
	_, err := s.db.Exec("INSERT INTO USERS (name) VALUES (?)", u.Name)
	return err
}

func (s *Store) GetUser(id int) (models.User, error) {
	var u models.User
	err := s.db.QueryRow("SELECT id, name FROM USERS WHERE id = ?", id).Scan(&u.ID, &u.Name)
	return u, err
}
