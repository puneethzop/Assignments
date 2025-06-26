package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	hellohandler(w, req)

	if w.Body.String() != "Hello, World!" {
		t.Error("Hello world failed")
	}
}

func openTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		t.Fatal("Failed to open DB:", err)
	}
	if err = db.Ping(); err != nil {
		t.Fatal("Failed to ping DB:", err)
	}
	return db
}

func TestSQL(t *testing.T) {
	db := &input{}
	db.data = openTestDB(t)
	defer db.data.Close()

	testPayload := `{ "task": "Testing"}`

	var initialCount int
	_ = db.data.QueryRow("SELECT count(*) FROM TASKS").Scan(&initialCount)

	req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(testPayload))
	resp := httptest.NewRecorder()
	db.addTask(resp, req)

	var newCount int
	_ = db.data.QueryRow("SELECT count(*) FROM TASKS").Scan(&newCount)
	if newCount != initialCount+1 {
		t.Errorf("Expected count %d, got %d", initialCount+1, newCount)
	}

	var lastID int
	_ = db.data.QueryRow("SELECT id FROM TASKS ORDER BY id DESC LIMIT 1").Scan(&lastID)

	// Test GET by ID
	req = httptest.NewRequest(http.MethodGet, "/task/{id}", nil)
	req.SetPathValue("id", strconv.Itoa(lastID))
	resp = httptest.NewRecorder()
	db.getByID(resp, req)

	expected := "ID: " + strconv.Itoa(lastID) + ", Task: Testing, Completed: false"
	if resp.Body.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, resp.Body.String())
	}

	// Test View Task (no effect check)
	req = httptest.NewRequest(http.MethodPut, "/task/{id}", nil)
	resp = httptest.NewRecorder()
	db.viewTask(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	// Test Mark as Complete
	req = httptest.NewRequest(http.MethodPut, "/task/{id}", nil)
	req.SetPathValue("id", strconv.Itoa(lastID))
	resp = httptest.NewRecorder()
	db.completeTask(resp, req)

	var completed bool
	_ = db.data.QueryRow("SELECT completed FROM TASKS WHERE id = ?", lastID).Scan(&completed)
	if !completed {
		t.Error("Task was not marked as completed")
	}

	// Test Delete
	req = httptest.NewRequest(http.MethodDelete, "/task/{id}", nil)
	req.SetPathValue("id", strconv.Itoa(lastID))
	resp = httptest.NewRecorder()
	db.deleteTask(resp, req)

	var idAfterDelete int
	err := db.data.QueryRow("SELECT id FROM TASKS ORDER BY id DESC LIMIT 1").Scan(&idAfterDelete)
	if err != nil {
		t.Errorf("Error fetching last ID after delete: %v", err)
	}
	if idAfterDelete == lastID {
		t.Errorf("Task not deleted: expected ID != %d, got %d", lastID, idAfterDelete)
	}
}

func TestSQLWithError(t *testing.T) {
	db := &input{}
	test := `{"Testing"}`

	//Opening DataBase
	var err error
	db.data, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.data.Close()
	err = db.data.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//ADD TASK when the JSON is corrupted
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/task", strings.NewReader(test))
	response := httptest.NewRecorder()

	db.addTask(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusBadRequest, response.Code)
	}

	//ADD TASK when the body is sent not right
	test = `{ "task": 1}`

	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/task", strings.NewReader(test))
	response = httptest.NewRecorder()

	db.addTask(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusBadRequest, response.Code)
	}

	//GET BY ID err check when ID is corrupted
	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", "r")

	db.getByID(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusBadRequest, response.Code)
	}

	//GET BY ID err check when ID is not present
	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", "0")

	db.getByID(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusNotFound, response.Code)
	}

	//COMPLETE TASK err check when ID is corrupted
	request = httptest.NewRequest(http.MethodPut, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", "r")

	db.completeTask(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusBadRequest, response.Code)
	}

	//COMPLETE TASK err check when ID is not present
	request = httptest.NewRequest(http.MethodPut, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", "0")

	db.completeTask(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusNotFound, response.Code)
	}

	//DELETE err check when ID is corrupted
	request = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", "r")

	db.deleteTask(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusBadRequest, response.Code)
	}

	//COMPLETE TASK err check when ID is not present
	request = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", "0")

	db.deleteTask(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Wrong result: expected %d, got %d", http.StatusNotFound, response.Code)
	}

}
