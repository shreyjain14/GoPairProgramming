package main

import "fmt"

func main() {

	var year int
	fmt.Print("Enter a year: ")
	fmt.Scan(&year)

	if year%4 == 0 {
		fmt.Println("It's a leap year.")
	} else {
		fmt.Println("It's not a leap year.")
	}

}
