package userhandler

import (
	"bytes"
	"encoding/json"
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

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		desc       string
		body       string
		mockErr    error
		wantStatus int
		wantBody   string
	}{
		{
			desc:       "valid user",
			body:       `{"id":1,"name":"John"}`,
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			desc:       "empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest,
			wantBody:   "Empty or unreadable",
		},
		{
			desc:       "invalid JSON",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid JSON",
		},
		{
			desc:       "service error",
			body:       `{"id":2,"name":"Error"}`,
			mockErr:    errors.New("insert failed"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "insert failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tc.body))
			w := httptest.NewRecorder()

			if tc.mockErr != nil || (tc.body != "" && tc.body != `{`) {
				var u models.User
				_ = json.Unmarshal([]byte(tc.body), &u)
				mockService.EXPECT().CreateUser(u).Return(tc.mockErr).AnyTimes()
			}

			handler.CreateUser(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantBody != "" && !strings.Contains(w.Body.String(), tc.wantBody) {
				t.Errorf("expected body to contain %q, got %q", tc.wantBody, w.Body.String())
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		desc       string
		id         string
		mockUser   models.User
		mockErr    error
		wantStatus int
		wantBody   string
	}{
		{
			desc:       "valid ID",
			id:         "1",
			mockUser:   models.User{ID: 1, Name: "John"},
			mockErr:    nil,
			wantStatus: http.StatusOK,
			wantBody:   `"id":1`,
		},
		{
			desc:       "missing ID",
			id:         "",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Missing ID",
		},
		{
			desc:       "invalid ID format",
			id:         "abc",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid ID format",
		},
		{
			desc:       "service error",
			id:         "2",
			mockErr:    errors.New("not found"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users/"+tc.id, nil)
			if tc.id != "" {
				req.SetPathValue("id", tc.id)
			}
			w := httptest.NewRecorder()

			if isNumber(tc.id) {
				id, _ := strconv.Atoi(tc.id)
				mockService.EXPECT().GetUser(id).Return(tc.mockUser, tc.mockErr).AnyTimes()
			}

			handler.GetUser(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantBody != "" && !strings.Contains(w.Body.String(), tc.wantBody) {
				t.Errorf("expected body to contain %q, got %q", tc.wantBody, w.Body.String())
			}
		})
	}
}
