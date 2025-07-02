// Make sure this file is named store_test.go
package userstore // Ensure this matches your user.go package if they are in the same directory

import (
	"3layerarch/models" // Adjust this import path if your models package is located elsewhere
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name     string
		user     models.User
		mockfunc func()
		err      error
	}{
		{
			name: "Successful CreateUser",
			user: models.User{Name: "TestUser"},
			mockfunc: func() {
				mock.SQL.ExpectExec("INSERT INTO USERS (name) VALUES (?)").
					WithArgs("TestUser").
					WillReturnResult(mock.SQL.NewResult(1, 1))
			},
			err: nil,
		},
		{
			name: "Failed CreateUser - Duplicate Entry",
			user: models.User{Name: "TestUser"},
			mockfunc: func() {
				mock.SQL.ExpectExec("INSERT INTO USERS (name) VALUES (?)").
					WithArgs("TestUser").
					WillReturnError(sql.ErrNoRows) // Simulating a database error
			},
			err: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			var db *sql.DB // db is not used directly, only through ctx.SQL
			// Assuming 'repo' is the variable name used in your store_test.go
			// If it's 'store', keep 'store'. The error message says 'repo'.
			repo := New(db)

			// FIX: Pass 'ctx' as the first argument
			err := repo.CreateUser(ctx, tt.user)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name     string
		userID   int
		mockfunc func()
		expAns   models.User
		err      error
	}{
		{
			name:   "Successful GetUser",
			userID: 1,
			mockfunc: func() {
				rows := mock.SQL.NewRows([]string{"id", "name"}).AddRow(1, "TestUser")
				mock.SQL.ExpectQuery("SELECT id, name FROM USERS WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expAns: models.User{ID: 1, Name: "TestUser"},
			err:    nil,
		},
		{
			name:   "GetUser - User Not Found",
			userID: 99,
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, name FROM USERS WHERE id = ?").
					WithArgs(99).
					WillReturnError(sql.ErrNoRows)
			},
			expAns: models.User{},
			err:    sql.ErrNoRows,
		},
		{
			name:   "Failed GetUser - Database Error",
			userID: 1,
			mockfunc: func() {
				mock.SQL.ExpectQuery("SELECT id, name FROM USERS WHERE id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone) // Simulating a connection error
			},
			expAns: models.User{},
			err:    sql.ErrConnDone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfunc()

			var db *sql.DB // db is not used directly, only through ctx.SQL
			// Assuming 'repo' is the variable name used in your store_test.go
			// If it's 'store', keep 'store'. The error message says 'repo'.
			repo := New(db)

			// FIX: Pass 'ctx' as the first argument
			ans, err := repo.GetUser(ctx, tt.userID)
			if !assert.Equal(t, tt.err, err) {
				t.Errorf("%v : error = %v, wantErr %v", tt.name, err, tt.err)
			}

			if !assert.Equal(t, tt.expAns, ans) {
				t.Errorf("%v : \nExpected = %v\n got = %v", tt.name, tt.expAns, ans)
			}
		})
	}
}
