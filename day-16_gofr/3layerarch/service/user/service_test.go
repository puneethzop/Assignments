package userservice

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"3layerarch/models"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := NewMockUserStore(ctrl)
	svc := New(mockStore)

	tests := []struct {
		desc      string
		input     models.User
		setupMock func()
		wantErr   error
	}{
		{
			desc:  "Success",
			input: models.User{Name: "Alice"},
			setupMock: func() {
				mockStore.EXPECT().CreateUser(models.User{Name: "Alice"}).Return(nil)
			},
			wantErr: nil,
		},
		{
			desc:    "Empty Name",
			input:   models.User{Name: ""},
			wantErr: errors.New("user name cannot be empty"),
		},
	}

	for _, test := range tests {
		if test.setupMock != nil {
			test.setupMock()
		}
		err := svc.CreateUser(test.input)
		if !errorsEqual(err, test.wantErr) {
			t.Errorf("%s: expected error %v, got %v", test.desc, test.wantErr, err)
		}
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := NewMockUserStore(ctrl)
	svc := New(mockStore)

	tests := []struct {
		desc    string
		inputID int
		setup   func()
		want    models.User
		wantErr error
	}{
		{
			desc:    "Success",
			inputID: 1,
			setup: func() {
				mockStore.EXPECT().GetUser(1).Return(models.User{ID: 1, Name: "Alice"}, nil)
			},
			want: models.User{ID: 1, Name: "Alice"},
		},
		{
			desc:    "Invalid ID",
			inputID: 0,
			wantErr: errors.New("invalid user ID"),
		},
		{
			desc:    "User Not Found",
			inputID: 2,
			setup: func() {
				mockStore.EXPECT().GetUser(2).Return(models.User{}, sql.ErrNoRows)
			},
			wantErr: errors.New("user not found"),
		},
		{
			desc:    "DB Error",
			inputID: 3,
			setup: func() {
				mockStore.EXPECT().GetUser(3).Return(models.User{}, errors.New("some db error"))
			},
			wantErr: errors.New("some db error"),
		},
	}

	for _, test := range tests {
		if test.setup != nil {
			test.setup()
		}

		got, err := svc.GetUser(test.inputID)

		if !errorsEqual(err, test.wantErr) {
			t.Errorf("%s: expected error %v, got %v", test.desc, test.wantErr, err)
		}

		if err == nil && !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: expected user %+v, got %+v", test.desc, test.want, got)
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
