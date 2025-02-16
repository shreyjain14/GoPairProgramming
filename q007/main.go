package main

import "fmt"

func main() {

	var n int
	fmt.Scan(&n)
	fmt.Println(factorial(n))

}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}
