package main

import "fmt"

func main() {

	var a, b, c int
	fmt.Print("Enter three numbers: ")
	fmt.Scan(&a, &b, &c)
	fmt.Printf("You entered: %d, %d, %d\n", a, b, c)

	if a > b && a > c {
		fmt.Printf("%d is the largest number\n", a)
	} else if b > a && b > c {
		fmt.Printf("%d is the largest number\n", b)
	} else {
		fmt.Printf("%d is the largest number\n", c)
	}

}
