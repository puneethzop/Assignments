package main

import (
  "fmt"
  "os"
  "strings"
  "time"
)

func main() {
  // Check if filename is passed as argument
  if len(os.Args) < 2 {
    fmt.Println("Usage: go run main.go <logfile.txt>")
    return
  }

  filename := os.Args[1]

  // Read entire file into memory
  data, err := os.ReadFile(filename)
  if err != nil {
    fmt.Println("Error reading file:", err)
    return
  }

  // Split into lines
  lines := strings.Split(string(data), "\n")

  // Initialize counters
  infoCount := 0
  warnCount := 0
  errCount := 0

  for _, line := range lines {
    switch {
    case strings.Contains(line, "[INFO]"):
      infoCount++
    case strings.Contains(line, "[WARN]"):
      warnCount++
    case strings.Contains(line, "[ERROR]"):
      errCount++
    }
  }

  // Printing summary
  fmt.Printf("Log Analysis of file : %s\n\n", filename)
  fmt.Println("Total number of  lines :", len(lines))
  fmt.Println("INFO :", infoCount, "entries")
  fmt.Println("WARNING :", warnCount, "entries")
  fmt.Println("ERROR :", errCount, "entries")

  // Printing timestamp
  fmt.Println("\nAnalyzed at:", time.Now().Format("2006-01-02 15:04:05"))
  fmt.Println("Thank you :)")
}
