package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Logger interface
type Logger interface {
	Log(message string)
}

// ConsoleLogger logs to console
type ConsoleLogger struct{}

func (c ConsoleLogger) Log(message string) {
	fmt.Println("Console:", message)
}

// FileLogger simulates file logging using a slice
type FileLogger struct {
	Logs []string
}

func (f *FileLogger) Log(message string) {
	f.Logs = append(f.Logs, message)
	fmt.Println("File:", message)
}

// RemoteLogger simulates remote logging
type RemoteLogger struct{}

func (r RemoteLogger) Log(message string) {
	fmt.Println("Remote:", message)
}

// LogAll sends message to all loggers
func LogAll(loggers []Logger, message string) {
	for _, logger := range loggers {
		logger.Log(message)
	}
}

func main() {
	// Setup loggers
	consoleLogger := ConsoleLogger{}
	fileLogger := &FileLogger{}
	remoteLogger := RemoteLogger{}

	loggers := []Logger{consoleLogger, fileLogger, remoteLogger}

	// Take user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a message to log: ")
	input, _ := reader.ReadString('\n')
	message := strings.TrimSpace(input)

	// Log to all
	LogAll(loggers, message)
}
