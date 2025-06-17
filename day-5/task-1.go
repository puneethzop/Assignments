package main

import (
	"fmt"
	"math"
)

// Shape interface
type Shape interface {
	Area() float64
}

// Circle struct
type Circle struct {
	Radius float64
}

// Rectangle struct
type Rectangle struct {
	Width  float64
	Height float64
}

// Area method for Circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Area method for Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Function that takes a Shape and prints its area
func PrintArea(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
}

func main() {
	var radius, width, height float64

	// Circle input
	fmt.Print("Enter radius of the circle: ")
	fmt.Scan(&radius)
	c := Circle{Radius: radius}

	// Rectangle input
	fmt.Print("Enter width of the rectangle: ")
	fmt.Scan(&width)
	fmt.Print("Enter height of the rectangle: ")
	fmt.Scan(&height)
	r := Rectangle{Width: width, Height: height}

	// Print Areas
	fmt.Print("Area of Circle: ")
	PrintArea(c)

	fmt.Print("Area of Rectangle: ")
	PrintArea(r)
}
