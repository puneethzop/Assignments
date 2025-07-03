package taskstore

import (
	"3layerarch/models"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"testing"
)

func TestStore_CreateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name     string
		task     models.Task
		mockfunc func()
		err      error
	}{
		{
			name: "Successful CreateTask",
			task: models.Task{Task: "Buy groceries", Completed: false, UserID: 1},
			mockfunc: func() {
				mock.SQL.ExpectExec("INSERT INTO TASKS (task, completed, user_id) VALUES (?, ?, ?)").
					WithArgs("Buy groceries", false, 1).
					WillReturnResult(mock.SQL.NewResult(1, 1))
			},
			err: nil,
		},
		{
			name: "Failed CreateTask - Duplicate entry (example)",
			task: models.Task{Task: "Buy groceries", Completed: false, UserID: 1},
			mockfunc: func() {
				mock.SQL.ExpectExec("INSERT INTO TASKS (task, completed, user_id) VALUES (?, ?, ?)").
					WithArgs("Buy groceries", false, 1).
					WillReturnError(errors.New("duplicate entry"))
			},
			err: errors.New("duplicate entry"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			svc := New()

			err := svc.CreateTask(ctx, tt.task)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestStore_GetTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	rows := mock.SQL.NewRows([]string{"id", "task", "completed", "user_id"}).
		AddRow(1, "Read a book", false, 1)

	tests := []struct {
		name     string
		taskID   int
		mockfunc func()
		expAns   models.Task
		err      error
	}{
		{
			name:   "Successful GetTask",
			taskID: 1,
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expAns: models.Task{ID: 1, Task: "Read a book", Completed: false, UserID: 1},
			err:    nil,
		},
		{
			name:   "Failed GetTask - No rows",
			taskID: 99,
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS WHERE id = ?").
					WithArgs(99).
					WillReturnError(sql.ErrNoRows) // Simulate no rows found
			},
			expAns: models.Task{},
			err:    sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			svc := New()

			ans, err := svc.GetTask(ctx, tt.taskID)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}

			if !assert.Equal(t, tt.expAns, ans) {
				t.Errorf("%v : \nExpected = %v\n got = %v", tt.name, tt.expAns, ans)
			}
		})
	}
}

func TestStore_ViewTasks(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	rows := mock.SQL.NewRows([]string{"id", "task", "completed", "user_id"}).
		AddRow(1, "Task 1", false, 1).
		AddRow(2, "Task 2", true, 2)

	tests := []struct {
		name     string
		mockfunc func()
		expAns   []models.Task
		err      error
	}{
		{
			name: "Successful ViewTasks",
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS").
					WillReturnRows(rows)
			},
			expAns: []models.Task{
				{ID: 1, Task: "Task 1", Completed: false, UserID: 1},
				{ID: 2, Task: "Task 2", Completed: true, UserID: 2},
			},
			err: nil,
		},
		{
			name: "Failed ViewTasks - DB error",
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS").
					WillReturnError(errors.New("db connection error"))
			},
			expAns: nil,
			err:    errors.New("db connection error"),
		},
		{
			name: "Empty ViewTasks",
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS").
					WillReturnRows(mock.SQL.NewRows([]string{"id", "task", "completed", "user_id"}))
			},
			expAns: []models.Task{},
			err:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			svc := New()

			ans, err := svc.ViewTasks(ctx)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}

			if !assert.Equal(t, tt.expAns, ans) {
				t.Errorf("%v : \nExpected = %v\n got = %v", tt.name, tt.expAns, ans)
			}
		})
	}
}
func TestStore_UpdateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name     string
		taskID   int
		mockfunc func()
		err      error
	}{
		{
			name:   "Successful UpdateTask",
			taskID: 1,
			mockfunc: func() {
				mock.SQL.ExpectExec("UPDATE TASKS SET completed = true WHERE id = ?").
					WithArgs(1).
					WillReturnResult(mock.SQL.NewResult(0, 1))
			},
			err: nil,
		},
		{
			name:   "Failed UpdateTask - No rows affected",
			taskID: 99,
			mockfunc: func() {
				mock.SQL.ExpectExec("UPDATE TASKS SET completed = true WHERE id = ?").
					WithArgs(99).
					WillReturnResult(mock.SQL.NewResult(0, 0))
			},
			err: nil,
		},
		{
			name:   "Failed UpdateTask - DB error",
			taskID: 1,
			mockfunc: func() {
				mock.SQL.ExpectExec("UPDATE TASKS SET completed = true WHERE id = ?").
					WithArgs(1).
					WillReturnError(errors.New("db write error"))
			},
			err: errors.New("db write error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			svc := New()

			err := svc.UpdateTask(ctx, tt.taskID)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestStore_DeleteTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name     string
		taskID   int
		mockfunc func()
		err      error
	}{
		{
			name:   "Successful DeleteTask",
			taskID: 1,
			mockfunc: func() {
				mock.SQL.ExpectExec("DELETE FROM TASKS WHERE id = ?").
					WithArgs(1).
					WillReturnResult(mock.SQL.NewResult(0, 1))
			},
			err: nil,
		},
		{
			name:   "Failed DeleteTask - No rows affected",
			taskID: 99,
			mockfunc: func() {
				mock.SQL.ExpectExec("DELETE FROM TASKS WHERE id = ?").
					WithArgs(99).
					WillReturnResult(mock.SQL.NewResult(0, 0))
			},
			err: nil,
		},
		{
			name:   "Failed DeleteTask - DB error",
			taskID: 1,
			mockfunc: func() {
				mock.SQL.ExpectExec("DELETE FROM TASKS WHERE id = ?").
					WithArgs(1).
					WillReturnError(errors.New("db delete error"))
			},
			err: errors.New("db delete error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			svc := New()

			err := svc.DeleteTask(ctx, tt.taskID)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}
		})
	}
}
