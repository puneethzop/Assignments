package main

import (
	"fmt"
)

// BankAccount struct
type BankAccount struct {
	Owner   string
	Balance float64
}

// Display balance (value receiver)
func (b BankAccount) DisplayBalance() {
	fmt.Printf("Owner: %s, Balance: %.2f\n", b.Owner, b.Balance)
}

// Deposit money (pointer receiver)
func (b *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("Deposit amount must be positive.")
		return
	}
	b.Balance += amount
	fmt.Printf("Deposited %.2f to %s's account\n", amount, b.Owner)
}

// Withdraw money (pointer receiver)
func (b *BankAccount) Withdraw(amount float64) {
	if amount > b.Balance {
		fmt.Println("Insufficient funds.")
		return
	}
	b.Balance -= amount
	fmt.Printf("Withdrew %.2f from %s's account\n", amount, b.Owner)
}

func main() {
	var account BankAccount
	account.Owner = "Puneeth"
	account.Balance = 1000.0

	var choice int

	for {
		fmt.Println("\n--- Bank Menu ---")
		fmt.Println("1. Display Balance")
		fmt.Println("2. Deposit")
		fmt.Println("3. Withdraw")
		fmt.Println("4. Exit")
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			account.DisplayBalance()
		case 2:
			var amount float64
			fmt.Print("Enter deposit amount: ")
			fmt.Scanln(&amount)
			account.Deposit(amount)
		case 3:
			var amount float64
			fmt.Print("Enter withdrawal amount: ")
			fmt.Scanln(&amount)
			account.Withdraw(amount)
		case 4:
			fmt.Println("Exiting. Thank you!")
			return
		default:
			fmt.Println("Invalid option. Please choose again.")
		}
	}
}
