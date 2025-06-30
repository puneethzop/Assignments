package taskservice_test

import (
	"database/sql"
	"errors"
	"testing"

	"3layerarch/models"
	"3layerarch/service/task"
)

// MockTaskStore implements TaskStore interface with function fields
type MockTaskStore struct {
	CreateTaskFn func(t models.Task) error
	GetTaskFn    func(id int) (models.Task, error)
	ViewTasksFn  func() ([]models.Task, error)
	UpdateTaskFn func(id int) error
	DeleteTaskFn func(id int) error
}

func (m *MockTaskStore) CreateTask(t models.Task) error {
	return m.CreateTaskFn(t)
}

func (m *MockTaskStore) GetTask(id int) (models.Task, error) {
	return m.GetTaskFn(id)
}

func (m *MockTaskStore) ViewTasks() ([]models.Task, error) {
	return m.ViewTasksFn()
}

func (m *MockTaskStore) UpdateTask(id int) error {
	return m.UpdateTaskFn(id)
}

func (m *MockTaskStore) DeleteTask(id int) error {
	return m.DeleteTaskFn(id)
}

// MockUserService implements UserService interface
type MockUserService struct {
	GetUserFn func(id int) (models.User, error)
}

func (m *MockUserService) GetUser(id int) (models.User, error) {
	return m.GetUserFn(id)
}

// ---- TESTS ----

func TestCreateTask_Success(t *testing.T) {
	mockStore := &MockTaskStore{
		CreateTaskFn: func(t models.Task) error {
			return nil
		},
	}
	mockUser := &MockUserService{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{ID: id, Name: "Alice"}, nil
		},
	}
	svc := taskservice.New(mockStore, mockUser)

	task := models.Task{
		Task:   "Test task",
		UserID: 1,
	}

	err := svc.CreateTask(task)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCreateTask_EmptyTask(t *testing.T) {
	svc := taskservice.New(nil, nil)

	task := models.Task{
		Task:   "",
		UserID: 1,
	}

	err := svc.CreateTask(task)
	if err == nil || err.Error() != "task cannot be empty" {
		t.Errorf("expected 'task cannot be empty' error, got %v", err)
	}
}

func TestCreateTask_UserNotFound(t *testing.T) {
	mockUser := &MockUserService{
		GetUserFn: func(id int) (models.User, error) {
			return models.User{}, errors.New("user not found")
		},
	}
	mockStore := &MockTaskStore{
		CreateTaskFn: func(t models.Task) error { return nil },
	}
	svc := taskservice.New(mockStore, mockUser)

	task := models.Task{
		Task:   "Valid task",
		UserID: 99,
	}

	err := svc.CreateTask(task)
	if err == nil || err.Error() != "user ID not found" {
		t.Errorf("expected 'user ID not found' error, got %v", err)
	}
}

func TestGetTask_Success(t *testing.T) {
	mockStore := &MockTaskStore{
		GetTaskFn: func(id int) (models.Task, error) {
			return models.Task{ID: id, Task: "task1", UserID: 1}, nil
		},
	}
	svc := taskservice.New(mockStore, nil)

	task, err := svc.GetTask(1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if task.ID != 1 {
		t.Errorf("expected task ID 1, got %d", task.ID)
	}
}

func TestGetTask_InvalidID(t *testing.T) {
	svc := taskservice.New(nil, nil)

	_, err := svc.GetTask(0)
	if err == nil || err.Error() != "invalid task ID" {
		t.Errorf("expected 'invalid task ID' error, got %v", err)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	mockStore := &MockTaskStore{
		GetTaskFn: func(id int) (models.Task, error) {
			return models.Task{}, sql.ErrNoRows
		},
	}
	svc := taskservice.New(mockStore, nil)

	_, err := svc.GetTask(1)
	if err == nil || err.Error() != "task not found" {
		t.Errorf("expected 'task not found' error, got %v", err)
	}
}

func TestViewTasks_Success(t *testing.T) {
	mockStore := &MockTaskStore{
		ViewTasksFn: func() ([]models.Task, error) {
			return []models.Task{
				{ID: 1, Task: "task1"},
				{ID: 2, Task: "task2"},
			}, nil
		},
	}
	svc := taskservice.New(mockStore, nil)

	tasks, err := svc.ViewTasks()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestViewTasks_Error(t *testing.T) {
	mockStore := &MockTaskStore{
		ViewTasksFn: func() ([]models.Task, error) {
			return nil, errors.New("db error")
		},
	}
	svc := taskservice.New(mockStore, nil)

	_, err := svc.ViewTasks()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected 'db error', got %v", err)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	mockStore := &MockTaskStore{
		GetTaskFn: func(id int) (models.Task, error) {
			return models.Task{ID: id}, nil
		},
		UpdateTaskFn: func(id int) error {
			return nil
		},
	}
	svc := taskservice.New(mockStore, nil)

	err := svc.UpdateTask(1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestUpdateTask_InvalidID(t *testing.T) {
	svc := taskservice.New(nil, nil)

	err := svc.UpdateTask(0)
	if err == nil || err.Error() != "invalid task ID" {
		t.Errorf("expected 'invalid task ID' error, got %v", err)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	mockStore := &MockTaskStore{
		GetTaskFn: func(id int) (models.Task, error) {
			return models.Task{}, sql.ErrNoRows // 🔧 proper sentinel error
		},
	}
	svc := taskservice.New(mockStore, nil)

	err := svc.UpdateTask(1)
	if err == nil || err.Error() != "task not found" {
		t.Errorf("expected 'task not found' error, got %v", err)
	}
}
func TestDeleteTask_Success(t *testing.T) {
	mockStore := &MockTaskStore{
		GetTaskFn: func(id int) (models.Task, error) {
			return models.Task{ID: id}, nil
		},
		DeleteTaskFn: func(id int) error {
			return nil
		},
	}
	svc := taskservice.New(mockStore, nil)

	err := svc.DeleteTask(1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestDeleteTask_InvalidID(t *testing.T) {
	svc := taskservice.New(nil, nil)

	err := svc.DeleteTask(0)
	if err == nil || err.Error() != "invalid task ID" {
		t.Errorf("expected 'invalid task ID' error, got %v", err)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockStore := &MockTaskStore{
		GetTaskFn: func(id int) (models.Task, error) {
			return models.Task{}, sql.ErrNoRows // Correct sentinel error
		},
	}
	svc := taskservice.New(mockStore, nil)

	err := svc.DeleteTask(1)
	if err == nil || err.Error() != "task not found" {
		t.Errorf("expected 'task not found' error, got %v", err)
	}
}
