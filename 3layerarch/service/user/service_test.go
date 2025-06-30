package userservice_test

import (
	"database/sql"
	"errors"
	"testing"

	"3layerarch/models"
	"3layerarch/service/user"
)

// MockUserStore implements UserStore interface with function fields
type MockUserStore struct {
	CreateUserFn func(u models.User) error
	GetUserFn    func(id int) (models.User, error)
}

func (m *MockUserStore) CreateUser(u models.User) error {
	return m.CreateUserFn(u)
}

func (m *MockUserStore) GetUser(id int) (models.User, error) {
	return m.GetUserFn(id)
}

func TestCreateUser_Success(t *testing.T) {
	mockStore := &MockUserStore{
		CreateUserFn: func(u models.User) error {
			return nil
		},
	}
	svc := userservice.New(mockStore)

	user := models.User{Name: "Alice"}

	err := svc.CreateUser(user)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	svc := userservice.New(nil)

	user := models.User{Name: ""}

	err := svc.CreateUser(user)
	if err == nil || err.Error() != "user name cannot be empty" {
		t.Errorf("expected 'user name cannot be empty' error, got %v", err)
	}
}

func TestGetUser_Success(t *testing.T) {
	mockStore := &MockUserStore{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{ID: id, Name: "Alice"}, nil
		},
	}
	svc := userservice.New(mockStore)

	user, err := svc.GetUser(1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if user.ID != 1 || user.Name != "Alice" {
		t.Errorf("unexpected user returned: %+v", user)
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	svc := userservice.New(nil)

	_, err := svc.GetUser(0)
	if err == nil || err.Error() != "invalid user ID" {
		t.Errorf("expected 'invalid user ID' error, got %v", err)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mockStore := &MockUserStore{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{}, sql.ErrNoRows
		},
	}
	svc := userservice.New(mockStore)

	_, err := svc.GetUser(1)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("expected 'user not found' error, got %v", err)
	}
}

func TestGetUser_OtherError(t *testing.T) {
	mockStore := &MockUserStore{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{}, errors.New("some db error")
		},
	}
	svc := userservice.New(mockStore)

	_, err := svc.GetUser(1)
	if err == nil || err.Error() != "some db error" {
		t.Errorf("expected 'some db error', got %v", err)
	}
}
