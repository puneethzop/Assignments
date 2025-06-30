package userstore_test

import (
	"3layerarch/models"
	"3layerarch/store/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	repo := userstore.New(db)
	user := models.User{Name: "Alice"}

	mock.ExpectExec("INSERT INTO USERS (name) VALUES (?)").
		WithArgs(user.Name).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	repo := userstore.New(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Bob")

	mock.ExpectQuery("SELECT id, name FROM USERS WHERE id = ?").
		WithArgs(1).WillReturnRows(rows)

	user, err := repo.GetUser(1)
	if err != nil || user.ID != 1 || user.Name != "Bob" {
		t.Errorf("unexpected result: %v, err: %v", user, err)
	}
}
