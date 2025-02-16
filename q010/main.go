package main

import "fmt"

func main() {
	var number int
	fmt.Scan(&number)

	sum := 0
	n := number

	if n < 0 {
		n = -n
	}

	for n > 0 {
		sum += n % 10
		n /= 10
	}

	fmt.Println(sum)
}
