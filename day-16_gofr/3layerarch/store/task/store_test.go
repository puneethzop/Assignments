package taskstore_test

import (
	"3layerarch/models"
	"3layerarch/store/task"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	defer db.Close()

	repo := taskstore.New(db)
	task := models.Task{Task: "Clean room", Completed: false, UserID: 1}

	mock.ExpectExec("INSERT INTO TASKS (task, completed, user_id) VALUES (?, ?, ?)").
		WithArgs(task.Task, task.Completed, task.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateTask(task)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	repo := taskstore.New(db)

	rows := sqlmock.NewRows([]string{"id", "task", "completed", "user_id"}).
		AddRow(1, "Read", false, 2)

	mock.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS WHERE id = ?").
		WithArgs(1).WillReturnRows(rows)

	task, err := repo.GetTask(1)
	if err != nil || task.ID != 1 {
		t.Errorf("unexpected result: %v, err: %v", task, err)
	}
}

func TestViewTasks(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	repo := taskstore.New(db)

	rows := sqlmock.NewRows([]string{"id", "task", "completed", "user_id"}).
		AddRow(1, "Work", false, 1)

	mock.ExpectQuery("SELECT id, task, completed, user_id FROM TASKS").
		WillReturnRows(rows)

	tasks, err := repo.ViewTasks()
	if err != nil || len(tasks) != 1 {
		t.Errorf("unexpected result: %v, err: %v", tasks, err)
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	repo := taskstore.New(db)

	mock.ExpectExec("UPDATE TASKS SET completed = true WHERE id = ?").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateTask(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	repo := taskstore.New(db)

	mock.ExpectExec("DELETE FROM TASKS WHERE id = ?").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteTask(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
