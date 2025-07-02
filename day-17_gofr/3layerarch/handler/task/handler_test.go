package taskhandler

import (
	"3layerarch/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http" // Alias to avoid conflict with net/http
)

// gofrResponse struct to hold expected results and errors
type gofrResponse struct {
	result any
	err    error
}

func TestCreateTask(t *testing.T) {
	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name        string
		requestBody string
		taskModel   models.Task // This struct should have a 'Task' field now
		expected    gofrResponse
		mockService bool // True if the service should be mocked
	}{
		{
			name:        "Success: Create a new task",
			requestBody: `{"id":1, "task":"Buy groceries", "completed":false}`,       // Changed "title" to "task"
			taskModel:   models.Task{ID: 1, Task: "Buy groceries", Completed: false}, // Changed Title to Task
			expected:    gofrResponse{result: "Task created", err: nil},
			mockService: true,
		},
		{
			name:        "Failure: Invalid JSON request body",
			requestBody: `{"id":1, "task": "Buy groceries", "completed":}`,            // Changed "title" to "task"
			taskModel:   models.Task{},                                                // Not used in this case as bind fails
			expected:    gofrResponse{result: nil, err: errors.New("unexpected EOF")}, // Example error from json.Unmarshal
			mockService: false,
		},
		{
			name:        "Failure: Service returns an error",
			requestBody: `{"id":2, "task":"Read a book", "completed":true}`,       // Changed "title" to "task"
			taskModel:   models.Task{ID: 2, Task: "Read a book", Completed: true}, // Changed Title to Task
			expected:    gofrResponse{result: nil, err: errors.New("database error")},
			mockService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockTaskService(ctrl)
			handler := New(mockService)

			req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			if tt.mockService {
				if tt.expected.err != nil && tt.expected.result == nil { // Expecting an error from service
					mockService.EXPECT().CreateTask(ctx, tt.taskModel).Return(tt.expected.err)
				} else { // Expecting success from service
					mockService.EXPECT().CreateTask(ctx, tt.taskModel).Return(nil)
				}
			}

			result, err := handler.CreateTask(ctx)

			assert.Equal(t, tt.expected.result, result, "Mismatch in result")
			if tt.expected.err != nil {
				// For Gofr bind errors, the error type might be specific.
				// For generic errors, just compare the error string.
				if _, ok := err.(*json.SyntaxError); ok { // Catch common JSON parsing errors
					assert.IsType(t, &json.SyntaxError{}, err, "Mismatch in error type")
				} else if _, ok := err.(*json.UnmarshalTypeError); ok {
					assert.IsType(t, &json.UnmarshalTypeError{}, err, "Mismatch in error type")
				} else {
					assert.EqualError(t, err, tt.expected.err.Error(), "Mismatch in error message")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name        string
		taskIDParam string
		taskID      int // Used for service mock expectation
		expected    gofrResponse
		mockService bool
	}{
		{
			name:        "Success: Get task by ID",
			taskIDParam: "1",
			taskID:      1,
			expected:    gofrResponse{result: models.Task{ID: 1, Task: "Test Task 1", Completed: false}, err: nil}, // Changed Title to Task
			mockService: true,
		},
		{
			name:        "Failure: Invalid ID parameter",
			taskIDParam: "abc",
			taskID:      0,                                                                                                   // Not relevant as strconv.Atoi fails
			expected:    gofrResponse{result: nil, err: &strconv.NumError{Func: "Atoi", Num: "abc", Err: strconv.ErrSyntax}}, // Specific error type
			mockService: false,
		},
		{
			name:        "Failure: Task not found (service returns error)",
			taskIDParam: "99",
			taskID:      99,
			expected:    gofrResponse{result: nil, err: errors.New("task not found")},
			mockService: true,
		},
		{
			name:        "Failure: Service returns generic error",
			taskIDParam: "3",
			taskID:      3,
			expected:    gofrResponse{result: nil, err: errors.New("database connection failed")},
			mockService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockTaskService(ctrl)
			handler := New(mockService)

			req := httptest.NewRequest(http.MethodGet, "/task/"+tt.taskIDParam, nil)

			// Use mux.SetURLVars to set path parameters on the underlying *http.Request
			vars := map[string]string{
				"id": tt.taskIDParam,
			}
			req = mux.SetURLVars(req, vars)

			ctx.Request = gofrHttp.NewRequest(req)

			if tt.mockService {
				// Conditionally set the return values for the mock based on whether an error is expected
				if tt.expected.err != nil {
					mockService.EXPECT().GetTask(ctx, tt.taskID).Return(models.Task{}, tt.expected.err)
				} else {
					mockService.EXPECT().GetTask(ctx, tt.taskID).Return(tt.expected.result.(models.Task), nil)
				}
			}

			result, err := handler.GetTask(ctx)

			assert.Equal(t, tt.expected.result, result, "Mismatch in result")
			if tt.expected.err != nil {
				// Assert specific error type for strconv.Atoi if needed
				if _, ok := tt.expected.err.(*strconv.NumError); ok {
					assert.IsType(t, tt.expected.err, err, "Mismatch in error type for Atoi")
					assert.Equal(t, tt.expected.err.Error(), err.Error(), "Mismatch in error message for Atoi")
				} else {
					assert.EqualError(t, err, tt.expected.err.Error(), "Mismatch in error message")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

func TestViewTasks(t *testing.T) {
	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name        string
		expected    gofrResponse
		mockService bool
	}{
		{
			name: "Success: View all tasks",
			expected: gofrResponse{
				result: []models.Task{
					{ID: 1, Task: "Task A", Completed: false}, // Changed Title to Task
					{ID: 2, Task: "Task B", Completed: true},  // Changed Title to Task
				},
				err: nil,
			},
			mockService: true,
		},
		{
			name:        "Success: No tasks available",
			expected:    gofrResponse{result: []models.Task{}, err: nil},
			mockService: true,
		},
		{
			name:        "Failure: Service returns an error",
			expected:    gofrResponse{result: nil, err: errors.New("failed to retrieve tasks")},
			mockService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockTaskService(ctrl)
			handler := New(mockService)

			req := httptest.NewRequest(http.MethodGet, "/task", nil)
			ctx.Request = gofrHttp.NewRequest(req)

			if tt.mockService {
				if tt.expected.err != nil {
					mockService.EXPECT().ViewTasks(ctx).Return(([]models.Task)(nil), tt.expected.err)
				} else {
					mockService.EXPECT().ViewTasks(ctx).Return(tt.expected.result.([]models.Task), nil)
				}
			}

			result, err := handler.ViewTasks(ctx)

			assert.Equal(t, tt.expected.result, result, "Mismatch in result")
			if tt.expected.err != nil {
				assert.EqualError(t, err, tt.expected.err.Error(), "Mismatch in error message")
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name        string
		taskIDParam string
		taskID      int
		expected    gofrResponse
		mockService bool
	}{
		{
			name:        "Success: Update task status",
			taskIDParam: "1",
			taskID:      1,
			expected:    gofrResponse{result: "Task updated", err: nil},
			mockService: true,
		},
		{
			name:        "Failure: Invalid ID parameter",
			taskIDParam: "xyz",
			taskID:      0,
			expected:    gofrResponse{result: nil, err: &strconv.NumError{Func: "Atoi", Num: "xyz", Err: strconv.ErrSyntax}},
			mockService: false,
		},
		{
			name:        "Failure: Task not found (service returns error)",
			taskIDParam: "99",
			taskID:      99,
			expected:    gofrResponse{result: nil, err: errors.New("task not found for update")},
			mockService: true,
		},
		{
			name:        "Failure: Service returns generic error",
			taskIDParam: "3",
			taskID:      3,
			expected:    gofrResponse{result: nil, err: errors.New("failed to update task in DB")},
			mockService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockTaskService(ctrl)
			handler := New(mockService)

			req := httptest.NewRequest(http.MethodPut, "/task/"+tt.taskIDParam, nil)
			// Use mux.SetURLVars to set path parameters on the underlying *http.Request
			vars := map[string]string{
				"id": tt.taskIDParam,
			}
			req = mux.SetURLVars(req, vars)

			ctx.Request = gofrHttp.NewRequest(req)

			if tt.mockService {
				mockService.EXPECT().UpdateTask(ctx, tt.taskID).Return(tt.expected.err)
			}

			result, err := handler.UpdateTask(ctx)

			assert.Equal(t, tt.expected.result, result, "Mismatch in result")
			if tt.expected.err != nil {
				if _, ok := tt.expected.err.(*strconv.NumError); ok {
					assert.IsType(t, tt.expected.err, err, "Mismatch in error type for Atoi")
					assert.Equal(t, tt.expected.err.Error(), err.Error(), "Mismatch in error message for Atoi")
				} else {
					assert.EqualError(t, err, tt.expected.err.Error(), "Mismatch in error message")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name        string
		taskIDParam string
		taskID      int
		expected    gofrResponse
		mockService bool
	}{
		{
			name:        "Success: Delete task",
			taskIDParam: "1",
			taskID:      1,
			expected:    gofrResponse{result: "Task deleted", err: nil},
			mockService: true,
		},
		{
			name:        "Failure: Invalid ID parameter",
			taskIDParam: "not-an-id",
			taskID:      0,
			expected:    gofrResponse{result: nil, err: &strconv.NumError{Func: "Atoi", Num: "not-an-id", Err: strconv.ErrSyntax}},
			mockService: false,
		},
		{
			name:        "Failure: Task not found (service returns error)",
			taskIDParam: "99",
			taskID:      99,
			expected:    gofrResponse{result: nil, err: errors.New("task not found for deletion")},
			mockService: true,
		},
		{
			name:        "Failure: Service returns generic error",
			taskIDParam: "3",
			taskID:      3,
			expected:    gofrResponse{result: nil, err: errors.New("failed to delete task from DB")},
			mockService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockTaskService(ctrl)
			handler := New(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/task/"+tt.taskIDParam, nil)
			// Use mux.SetURLVars to set path parameters on the underlying *http.Request
			vars := map[string]string{
				"id": tt.taskIDParam,
			}
			req = mux.SetURLVars(req, vars)

			ctx.Request = gofrHttp.NewRequest(req)

			if tt.mockService {
				mockService.EXPECT().DeleteTask(ctx, tt.taskID).Return(tt.expected.err)
			}

			result, err := handler.DeleteTask(ctx)

			assert.Equal(t, tt.expected.result, result, "Mismatch in result")
			if tt.expected.err != nil {
				if _, ok := tt.expected.err.(*strconv.NumError); ok {
					assert.IsType(t, tt.expected.err, err, "Mismatch in error type for Atoi")
					assert.Equal(t, tt.expected.err.Error(), err.Error(), "Mismatch in error message for Atoi")
				} else {
					assert.EqualError(t, err, tt.expected.err.Error(), "Mismatch in error message")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}
