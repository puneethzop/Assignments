package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	tm := &TaskManager{}

	// Test addTask
	body := `{"task":"wake up early"}`
	req := httptest.NewRequest("POST", "/task", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	tm.addTask(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("addTask: Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if len(tm.tasks) != 1 {
		t.Errorf("addTask: Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].Task != "wake up early" || tm.tasks[0].Completed {
		t.Errorf("addTask: Task was not added correctly: %+v", tm.tasks[0])
	}

	// Test getByID
	req = httptest.NewRequest("GET", "/task/0", nil)
	req.SetPathValue("id", "0")
	w = httptest.NewRecorder()
	tm.getByID(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("getByID: Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	expectedGet := `{"task":"wake up early","completed":false}`
	if strings.TrimSpace(w.Body.String()) != expectedGet {
		t.Errorf("getByID: Expected body %s, got %s", expectedGet, w.Body.String())
	}

	// Test viewAll
	req = httptest.NewRequest("GET", "/task", nil)
	w = httptest.NewRecorder()
	tm.viewAll(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("viewAll: Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	expectedAll := `[{"task":"wake up early","completed":false}]`
	if strings.TrimSpace(w.Body.String()) != expectedAll {
		t.Errorf("viewAll: Expected body %s, got %s", expectedAll, w.Body.String())
	}

	// Test completeTask
	req = httptest.NewRequest("PATCH", "/task/0", nil)
	req.SetPathValue("id", "0")
	w = httptest.NewRecorder()
	tm.completeTask(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("completeTask: Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if !tm.tasks[0].Completed {
		t.Errorf("completeTask: Task was not marked completed: %+v", tm.tasks[0])
	}

	// Test deleteTask
	req = httptest.NewRequest("DELETE", "/task/0", nil)
	req.SetPathValue("id", "0")
	w = httptest.NewRecorder()
	tm.deleteTask(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("deleteTask: Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if len(tm.tasks) != 0 {
		t.Errorf("deleteTask: Task was not deleted properly, tasks left: %d", len(tm.tasks))
	}

}
