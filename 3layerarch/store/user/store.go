package userstore

import (
	"3layerarch/models"
	"database/sql"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) Create(user models.User) error {
	_, err := s.DB.Exec("INSERT INTO USERS (name) VALUES (?)", user.Name)
	return err
}

func (s *Store) GetByID(id int) (models.User, error) {
	var user models.User
	err := s.DB.QueryRow("SELECT id, name FROM USERS WHERE id = ?", id).Scan(&user.ID, &user.Name)
	return user, err
}
