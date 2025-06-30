package userhandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"3layerarch/handler/user"
	"3layerarch/models"
)

// MockUserService implements UserService interface with function fields
type MockUserService struct {
	CreateUserFn func(u models.User) error
	GetUserFn    func(id int) (models.User, error)
}

func (m *MockUserService) CreateUser(u models.User) error {
	return m.CreateUserFn(u)
}

func (m *MockUserService) GetUser(id int) (models.User, error) {
	return m.GetUserFn(id)
}

func TestCreateUserHandler_Success(t *testing.T) {
	mockSvc := &MockUserService{
		CreateUserFn: func(u models.User) error { return nil },
	}
	handler := userhandler.New(mockSvc)

	user := models.User{Name: "Alice"}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", w.Result().StatusCode)
	}
}

func TestCreateUserHandler_EmptyBody(t *testing.T) {
	mockSvc := &MockUserService{}
	handler := userhandler.New(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %d", w.Result().StatusCode)
	}
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	mockSvc := &MockUserService{}
	handler := userhandler.New(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader([]byte(`{invalid json}`)))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %d", w.Result().StatusCode)
	}
}

func TestCreateUserHandler_ServiceError(t *testing.T) {
	mockSvc := &MockUserService{
		CreateUserFn: func(u models.User) error { return errors.New("service error") },
	}
	handler := userhandler.New(mockSvc)

	user := models.User{Name: "Alice"}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %d", w.Result().StatusCode)
	}
}

func TestGetUserHandler_Success(t *testing.T) {
	mockSvc := &MockUserService{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{ID: id, Name: "Alice"}, nil
		},
	}
	handler := userhandler.New(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	req.SetPathValue("id", "1") // This is the key fix
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", w.Result().StatusCode)
	}

	var u models.User
	err := json.NewDecoder(w.Body).Decode(&u)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if u.ID != 1 || u.Name != "Alice" {
		t.Errorf("unexpected user: %+v", u)
	}
}

func TestGetUserHandler_MissingID(t *testing.T) {
	mockSvc := &MockUserService{}
	handler := userhandler.New(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %d", w.Result().StatusCode)
	}
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	mockSvc := &MockUserService{}
	handler := userhandler.New(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user?id=abc", nil)
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %d", w.Result().StatusCode)
	}
}

func TestGetUserHandler_ServiceError(t *testing.T) {
	mockSvc := &MockUserService{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{}, errors.New("service error")
		},
	}
	handler := userhandler.New(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user?id=1", nil)
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %d", w.Result().StatusCode)
	}
}
