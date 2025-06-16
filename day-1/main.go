package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Hardcoded exchange rates
var exchangeRates = map[string]map[string]float64{
	"USD": {"USD": 1, "INR": 83.12, "EUR": 0.92, "JPY": 155.75},
	"INR": {"USD": 0.012, "INR": 1, "EUR": 0.011, "JPY": 1.87},
	"EUR": {"USD": 1.09, "INR": 90.32, "EUR": 1, "JPY": 169.80},
	"JPY": {"USD": 0.0064, "INR": 0.54, "EUR": 0.0059, "JPY": 1},
}

func main() {

	currentHour := time.Now().Hour()
	if currentHour < 12 {
		fmt.Println("Good morning! ")
	} else if currentHour < 18 {
		fmt.Println("Good afternoon! ")
	} else if currentHour < 21 {
		fmt.Println("Good evening! ")
	} else {
		fmt.Println("Good night! ")
	}

	if len(os.Args) != 4 {
		// If user didn't give 3 inputs, show message and stop
		fmt.Println("Please give amount, Source_currency and target_currency.")
		fmt.Println("Example: go run main.go 100 USD INR")
		return
	}

	// Parse inputs
	amountStr := os.Args[1]
	from := strings.ToUpper(os.Args[2])
	to := strings.ToUpper(os.Args[3])

	// Convert amount to float
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount < 0 {
		fmt.Println("Error: Amount must be a valid positive number:)")
		return
	}

	// Validate currencies
	if !isValidCurrency(from) || !isValidCurrency(to) {
		fmt.Println("Error: Supports only USD, INR, EUR, JPY Currencies")
		return
	}

	// Do conversion
	converted := convert(amount, from, to)
	fmt.Printf("%.2f %s is equivalent to %.2f %s\nThank You :)", amount, from, converted, to)
}

// will Check if currency exists in map
func isValidCurrency(currency string) bool {
	_, exists := exchangeRates[currency]
	return exists
}

// Converts from one currency to another
func convert(amount float64, from string, to string) float64 {
	rate := exchangeRates[from][to]
	return amount * rate
}
