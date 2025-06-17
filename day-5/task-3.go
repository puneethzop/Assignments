package main

import (
	"fmt"
	//"strings"
)

// PaymentMethod interface
type PaymentMethod interface {
	Pay(amount float64) string
}

// OTPEnabled interface (optional behavior)
type OTPEnabled interface {
	GenerateOTP()
}

// CreditCard struct
type CreditCard struct {
	CardNumber string
}

func (c CreditCard) Pay(amount float64) string {
	last4 := c.CardNumber[len(c.CardNumber)-4:]
	return fmt.Sprintf("[CreditCard] Paid ₹%.2f using card ending with %s", amount, last4)
}

func (c CreditCard) GenerateOTP() {
	fmt.Println("[CreditCard] OTP sent to registered number")
}

// PayPal struct
type PayPal struct {
	Email string
}

func (p PayPal) Pay(amount float64) string {
	return fmt.Sprintf("[PayPal] Paid ₹%.2f using PayPal account: %s", amount, p.Email)
}

// UPI struct
type UPI struct {
	UPIID string
}

func (u UPI) Pay(amount float64) string {
	return fmt.Sprintf("[UPI] Paid ₹%.2f using UPI: %s", amount, u.UPIID)
}

func (u UPI) GenerateOTP() {
	fmt.Println("[UPI] OTP sent to registered device")
}

// Main function
func main() {
	methods := []PaymentMethod{
		CreditCard{CardNumber: "1234567812341234"},
		PayPal{Email: "user@example.com"},
		UPI{UPIID: "user@upi"},
	}

	for _, method := range methods {
		// Check for optional OTP behavior
		if otpCapable, ok := method.(OTPEnabled); ok {
			otpCapable.GenerateOTP()
		}
		// Always call Pay
		fmt.Println(method.Pay(500))
		fmt.Println()
	}
}
