package main

import "fmt"

func main() {

	var a, b int
	var sign string

	fmt.Print("Enter first number: ")
	fmt.Scanln(&a)
	fmt.Print("Enter operator (+, -, *, /): ")
	fmt.Scanln(&sign)
	fmt.Print("Enter second number: ")
	fmt.Scanln(&b)

	switch sign {
	case "+":
		fmt.Printf("Result: %d\n", a+b)
	case "-":
		fmt.Printf("Result: %d\n", a-b)
	case "*":
		fmt.Printf("Result: %d\n", a*b)
	case "/":
		if b == 0 {
			fmt.Println("Error: Division by zero")
		} else {
			fmt.Printf("Result: %d\n", a/b)
		}
	default:
		fmt.Println("Invalid operator")
	}

}
