package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Task struct represents a task with ID, description and completion status
type Task struct {
	ID          int
	Description string
	Completed   bool
}

// Closure to generate unique IDs for every single task
func idGenerator() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

// AddTask adds a new task to the list
func AddTask(description string, tasks *[]*Task, getID func() int) {
	id := getID()
	newTask := &Task{
		ID:          id,
		Description: description,
		Completed:   false,
	}
	*tasks = append(*tasks, newTask)
	fmt.Printf("Task Added: %d - %s\n", id, description)
}

// ListTasks prints all pending (not completed) tasks
func ListTasks(tasks *[]*Task) {
	fmt.Println("\nPending Tasks:")
	for _, task := range *tasks {
		if !task.Completed {
			fmt.Printf("%d: %s\n", task.ID, task.Description)
		}
	}
}

// CompleteTask marks a task as completed based on its ID
func CompleteTask(id int, tasks *[]*Task) {
	for _, task := range *tasks {
		if task.ID == id {
			task.Completed = true
			fmt.Printf("Task %d marked as completed.\n", id)
			return
		}
	}
	fmt.Printf("Task with ID %d not found.\n", id)
}

func main() {
	tasks := []*Task{}
	getNextID := idGenerator()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nChoose an action: add, list, complete, exit")
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			fmt.Println("No input detected.\nPlease choose a command: add, list, complete, or exit.")
			continue
		}

		switch args[0] {
		case "add":
			if len(args) < 2 {
				fmt.Println("Usage: add <task description>")
				continue
			}
			description := strings.Join(args[1:], " ")
			AddTask(description, &tasks, getNextID)

		case "list":
			ListTasks(&tasks)

		case "complete":
			if len(args) != 2 {
				fmt.Println("Usage: complete <task ID>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid task ID.")
				continue
			}
			CompleteTask(id, &tasks)

		case "exit":
			fmt.Println("Exiting program...")
			fmt.Println("Thank you for using the Task Tracker!")
			return

		default:
			fmt.Println("Unknown command.")
		}
	}
}
