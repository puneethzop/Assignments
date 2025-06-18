package main

import (
	"testing"
)

func TestAddAndCompleteTask(t *testing.T) {
	// Test AddTask
	tasks := []*Task{}
	getNextID := idGenerator()

	AddTask("Write unit tests", &tasks, getNextID)

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Description != "Write unit tests" {
		t.Errorf("Expected 'Write unit tests', got '%s'", tasks[0].Description)
	}
	if tasks[0].Completed {
		t.Error("Expected task to be incomplete")
	}

	// Test CompleteTask
	CompleteTask(1, &tasks)

	if !tasks[0].Completed {
		t.Error("Expected task to be completed")
	}
}
