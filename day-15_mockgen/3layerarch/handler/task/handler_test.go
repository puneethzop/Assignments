package taskhandler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"3layerarch/models"
	"go.uber.org/mock/gomock"
)

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockTaskService(ctrl)
	handler := New(mockService)

	// valid id
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		mockService.EXPECT().GetTask(1).Return(models.Task{ID: 1, Task: "test"}, nil)

		handler.GetTask(w, req)
		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", res.StatusCode)
		}
		if !strings.Contains(w.Body.String(), `"id":1`) {
			t.Errorf("expected body to contain id:1, got %s", w.Body.String())
		}
	}

	// missing id
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.GetTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// invalid id
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()

		handler.GetTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks/2", nil)
		req.SetPathValue("id", "2")
		w := httptest.NewRecorder()

		mockService.EXPECT().GetTask(2).Return(models.Task{}, errors.New("task not found"))

		handler.GetTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}
}

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockTaskService(ctrl)
	handler := New(mockService)

	// valid
	{
		body := `{"task":"new task","completed":false,"user_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		task := models.Task{Task: "new task", Completed: false, UserID: 1}
		mockService.EXPECT().CreateTask(task).Return(nil)

		handler.CreateTask(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	}

	// empty body
	{
		req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// invalid json
	{
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString("{"))
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// service error
	{
		body := `{"task":"fail","completed":false,"user_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		task := models.Task{Task: "fail", Completed: false, UserID: 1}
		mockService.EXPECT().CreateTask(task).Return(errors.New("fail"))

		handler.CreateTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}
}

func TestViewTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockTaskService(ctrl)
	handler := New(mockService)

	// success
	{
		mockService.EXPECT().ViewTasks().Return([]models.Task{{ID: 1, Task: "t"}}, nil)
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.ViewTasks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), `"id":1`) {
			t.Errorf("expected body to contain task id 1")
		}
	}

	// failure
	{
		mockService.EXPECT().ViewTasks().Return(nil, errors.New("db error"))
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.ViewTasks(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockTaskService(ctrl)
	handler := New(mockService)

	// valid
	{
		req := httptest.NewRequest(http.MethodPut, "/tasks/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		mockService.EXPECT().UpdateTask(1).Return(nil)

		handler.UpdateTask(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	}

	// missing id
	{
		req := httptest.NewRequest(http.MethodPut, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.UpdateTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// invalid id
	{
		req := httptest.NewRequest(http.MethodPut, "/tasks/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()

		handler.UpdateTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodPut, "/tasks/2", nil)
		req.SetPathValue("id", "2")
		w := httptest.NewRecorder()

		mockService.EXPECT().UpdateTask(2).Return(errors.New("not found"))

		handler.UpdateTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockTaskService(ctrl)
	handler := New(mockService)

	// valid
	{
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		mockService.EXPECT().DeleteTask(1).Return(nil)

		handler.DeleteTask(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	}

	// missing id
	{
		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.DeleteTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// invalid id
	{
		req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()

		handler.DeleteTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodDelete, "/tasks/2", nil)
		req.SetPathValue("id", "2")
		w := httptest.NewRecorder()

		mockService.EXPECT().DeleteTask(2).Return(errors.New("not found"))

		handler.DeleteTask(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	}
}
