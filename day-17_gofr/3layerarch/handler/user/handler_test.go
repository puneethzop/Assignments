package userhandler

import (
	"3layerarch/models"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name        string
		requestBody string
		inputUser   models.User
		expected    gofrResponse
		ifMock      bool
	}{
		{
			name:        "Success",
			requestBody: `{"id":1,"name":"Tester"}`,
			inputUser:   models.User{ID: 1, Name: "Tester"},
			expected: gofrResponse{
				result: "error",
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:        "Bind Failure",
			requestBody: `invalid-json`,
			expected: gofrResponse{
				result: nil,
				err:    nil, // We'll check this using ErrorContains
			},
			ifMock: false,
		},
		{
			name:        "Service Error",
			requestBody: `{"id":1,"name":"Tester"}`,
			inputUser:   models.User{ID: 1, Name: "Tester"},
			expected: gofrResponse{
				result: nil,
				err:    errors.New("service error"),
			},
			ifMock: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := NewMockUserService(ctrl)
			h := New(mockService)

			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = gofrHttp.NewRequest(req)

			if tt.ifMock {
				mockService.EXPECT().CreateUser(ctx, tt.inputUser).Return(tt.expected.err)
			}

			val, err := h.CreateUser(ctx)
			resp := gofrResponse{val, err}

			if tt.name == "Bind Failure" {
				assert.Nil(t, resp.result)
				assert.ErrorContains(t, err, "invalid character")
			} else {
				assert.Equal(t, tt.expected, resp, "Test[%d] Failed: %s", i, tt.name)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	tests := []struct {
		name     string
		pathID   string
		expected gofrResponse
		ifMock   bool
	}{
		{
			name:   "Success",
			pathID: "1",
			expected: gofrResponse{
				result: models.User{ID: 1, Name: "Tester"},
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:   "Invalid ID (non-integer)",
			pathID: "abc",
			expected: gofrResponse{
				result: nil,
				err:    nil, // We'll check error type in test
			},
			ifMock: false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := NewMockUserService(ctrl)
			h := New(mockService)

			req := httptest.NewRequest(http.MethodGet, "/user/{id}", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tt.pathID})
			ctx.Request = gofrHttp.NewRequest(req)

			if tt.ifMock {
				id, _ := strconv.Atoi(tt.pathID)
				mockService.EXPECT().GetUser(ctx, id).Return(tt.expected.result, tt.expected.err)
			}

			val, err := h.GetUser(ctx)
			resp := gofrResponse{val, err}

			if tt.name == "Invalid ID (non-integer)" {
				assert.Nil(t, resp.result)
				assert.ErrorContains(t, err, "invalid syntax")
			} else {
				assert.Equal(t, tt.expected, resp, "Test[%d] Failed: %s", i, tt.name)
			}
		})
	}
}
