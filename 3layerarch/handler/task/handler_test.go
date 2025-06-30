package taskhandler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"3layerarch/handler/task"
	"3layerarch/models"
)

type MockService struct{}

func (m *MockService) CreateTask(t models.Task) error {
	if t.Task == "" {
		return io.EOF
	}
	return nil
}

func (m *MockService) GetTask(id int) (models.Task, error) {
	if id == 1 {
		return models.Task{ID: 1, Task: "Hello", Completed: false, UserID: 1}, nil
	}
	return models.Task{}, io.EOF
}

func (m *MockService) ViewTasks() ([]models.Task, error) {
	return []models.Task{
		{ID: 1, Task: "Test 1", Completed: false, UserID: 1},
		{ID: 2, Task: "Test 2", Completed: true, UserID: 2},
	}, nil
}

func (m *MockService) UpdateTask(id int) error {
	if id == 0 {
		return io.EOF
	}
	return nil
}

func (m *MockService) DeleteTask(id int) error {
	if id == 0 {
		return io.EOF
	}
	return nil
}

func TestCreateTaskHandler(t *testing.T) {
	handler := taskhandler.New(&MockService{})

	// Valid Task
	task := models.Task{Task: "Test task", UserID: 1}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.CreateTask(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Invalid Task (empty body)
	req = httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte{}))
	w = httptest.NewRecorder()
	handler.CreateTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty body, got %d", w.Code)
	}

	// Invalid Task (invalid JSON)
	req = httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte("{invalid json")))
	w = httptest.NewRecorder()
	handler.CreateTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid json, got %d", w.Code)
	}
}

func TestGetTaskHandler(t *testing.T) {
	handler := taskhandler.New(&MockService{})

	// Valid ID
	req := httptest.NewRequest(http.MethodGet, "/task/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	handler.GetTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Hello") {
		t.Errorf("unexpected response: %s", w.Body.String())
	}

	// Missing ID
	req = httptest.NewRequest(http.MethodGet, "/task/", nil)
	req.SetPathValue("id", "")
	w = httptest.NewRecorder()
	handler.GetTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing ID, got %d", w.Code)
	}

	// Invalid ID format
	req = httptest.NewRequest(http.MethodGet, "/task/abc", nil) // Fix here
	req.SetPathValue("id", "abc")                               // Fix here
	w = httptest.NewRecorder()
	handler.GetTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID format, got %d", w.Code)
	}

	// Not Found
	req = httptest.NewRequest(http.MethodGet, "/task/999", nil)
	req.SetPathValue("id", "999")
	w = httptest.NewRecorder()
	handler.GetTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for not found, got %d", w.Code)
	}
}

func TestViewTasksHandler(t *testing.T) {
	handler := taskhandler.New(&MockService{})

	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	w := httptest.NewRecorder()
	handler.ViewTasks(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Test 1") || !strings.Contains(w.Body.String(), "Test 2") {
		t.Errorf("unexpected tasks list response: %s", w.Body.String())
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	handler := taskhandler.New(&MockService{})

	// Valid ID
	req := httptest.NewRequest(http.MethodPut, "/task/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	handler.UpdateTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	// Missing ID
	req = httptest.NewRequest(http.MethodPut, "/task/", nil)
	req.SetPathValue("id", "")
	w = httptest.NewRecorder()
	handler.UpdateTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing ID, got %d", w.Code)
	}

	// Invalid ID
	req = httptest.NewRequest(http.MethodPut, "/task/abc", nil)
	req.SetPathValue("id", "abc")
	w = httptest.NewRecorder()
	handler.UpdateTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID, got %d", w.Code)
	}

	// Error in service for update
	req = httptest.NewRequest(http.MethodPut, "/task/0", nil)
	req.SetPathValue("id", "0")
	w = httptest.NewRecorder()
	handler.UpdateTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for update error, got %d", w.Code)
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	handler := taskhandler.New(&MockService{})

	// Valid ID
	req := httptest.NewRequest(http.MethodDelete, "/task/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	// Missing ID
	req = httptest.NewRequest(http.MethodDelete, "/task/", nil)
	req.SetPathValue("id", "")
	w = httptest.NewRecorder()
	handler.DeleteTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing ID, got %d", w.Code)
	}

	// Invalid ID
	req = httptest.NewRequest(http.MethodDelete, "/task/abc", nil)
	req.SetPathValue("id", "abc")
	w = httptest.NewRecorder()
	handler.DeleteTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID, got %d", w.Code)
	}

	// Error in service for delete
	req = httptest.NewRequest(http.MethodDelete, "/task/0", nil)
	req.SetPathValue("id", "0")
	w = httptest.NewRecorder()
	handler.DeleteTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for delete error, got %d", w.Code)
	}
}
