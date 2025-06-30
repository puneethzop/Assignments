package taskservice

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"3layerarch/models"
	"go.uber.org/mock/gomock"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTaskStore := NewMockTaskStore(ctrl)
	mockUserService := NewMockUserService(ctrl)
	svc := New(mockTaskStore, mockUserService)

	tests := []struct {
		desc      string
		task      models.Task
		ifMock    bool
		userFound bool
		storeErr  error
		userErr   error
		wantErr   error
	}{
		{
			"Valid Task", models.Task{Task: "Write tests", UserID: 1},
			true, true, nil, nil, nil,
		},
		{
			"Empty Task", models.Task{Task: "", UserID: 1},
			false, false, nil, nil, errors.New("task cannot be empty"),
		},
		{
			"User Not Found", models.Task{Task: "Write tests", UserID: 42},
			true, false, nil, errors.New("user not found"), errors.New("user ID not found"),
		},
	}

	for _, test := range tests {
		if test.ifMock {
			if test.userErr != nil {
				mockUserService.EXPECT().GetUser(test.task.UserID).Return(models.User{}, test.userErr)
			} else {
				mockUserService.EXPECT().GetUser(test.task.UserID).Return(models.User{ID: test.task.UserID}, nil)
				mockTaskStore.EXPECT().CreateTask(test.task).Return(test.storeErr)
			}
		}

		err := svc.CreateTask(test.task)
		if !errorsEqual(err, test.wantErr) {
			t.Errorf("%v: expected error %v, got %v", test.desc, test.wantErr, err)
		}
	}
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTaskStore := NewMockTaskStore(ctrl)
	svc := New(mockTaskStore, nil)

	tests := []struct {
		desc    string
		taskID  int
		task    models.Task
		err     error
		want    models.Task
		wantErr error
	}{
		{
			"Valid Task", 1, models.Task{ID: 1, Task: "Test", UserID: 1},
			nil, models.Task{ID: 1, Task: "Test", UserID: 1}, nil,
		},
		{
			"Invalid ID", 0, models.Task{}, nil, models.Task{}, errors.New("invalid task ID"),
		},
		{
			"Not Found", 2, models.Task{}, sql.ErrNoRows, models.Task{}, errors.New("task not found"),
		},
		{
			"DB Error", 3, models.Task{}, errors.New("db error"), models.Task{}, errors.New("db error"),
		},
	}

	for _, test := range tests {
		if test.taskID > 0 {
			mockTaskStore.EXPECT().GetTask(test.taskID).Return(test.task, test.err)
		}

		got, err := svc.GetTask(test.taskID)
		if !errorsEqual(err, test.wantErr) {
			t.Errorf("%v: expected error %v, got %v", test.desc, test.wantErr, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%v: expected task %v, got %v", test.desc, test.want, got)
		}
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTaskStore := NewMockTaskStore(ctrl)
	svc := New(mockTaskStore, nil)

	tests := []struct {
		desc      string
		taskID    int
		getErr    error
		updateErr error
		wantErr   error
	}{
		{
			"Valid Update", 1, nil, nil, nil,
		},
		{
			"Invalid ID", 0, nil, nil, errors.New("invalid task ID"),
		},
		{
			"Task Not Found", 2, sql.ErrNoRows, nil, errors.New("task not found"),
		},
		{
			"Get DB Error", 3, errors.New("db error"), nil, errors.New("db error"),
		},
		{
			"Update DB Error", 4, nil, errors.New("update failed"), errors.New("update failed"),
		},
	}

	for _, test := range tests {
		if test.taskID > 0 {
			mockTaskStore.EXPECT().GetTask(test.taskID).Return(models.Task{}, test.getErr)
			if test.getErr == nil {
				mockTaskStore.EXPECT().UpdateTask(test.taskID).Return(test.updateErr)
			}
		}

		err := svc.UpdateTask(test.taskID)
		if !errorsEqual(err, test.wantErr) {
			t.Errorf("%v: expected error %v, got %v", test.desc, test.wantErr, err)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTaskStore := NewMockTaskStore(ctrl)
	svc := New(mockTaskStore, nil)

	tests := []struct {
		desc      string
		taskID    int
		getErr    error
		deleteErr error
		wantErr   error
	}{
		{
			"Valid Delete", 1, nil, nil, nil,
		},
		{
			"Invalid ID", 0, nil, nil, errors.New("invalid task ID"),
		},
		{
			"Task Not Found", 2, sql.ErrNoRows, nil, errors.New("task not found"),
		},
		{
			"Get DB Error", 3, errors.New("db error"), nil, errors.New("db error"),
		},
		{
			"Delete DB Error", 4, nil, errors.New("delete failed"), errors.New("delete failed"),
		},
	}

	for _, test := range tests {
		if test.taskID > 0 {
			mockTaskStore.EXPECT().GetTask(test.taskID).Return(models.Task{}, test.getErr)
			if test.getErr == nil {
				mockTaskStore.EXPECT().DeleteTask(test.taskID).Return(test.deleteErr)
			}
		}

		err := svc.DeleteTask(test.taskID)
		if !errorsEqual(err, test.wantErr) {
			t.Errorf("%v: expected error %v, got %v", test.desc, test.wantErr, err)
		}
	}
}

func errorsEqual(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}
	return err1.Error() == err2.Error()
}
