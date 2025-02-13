package main

import "fmt"

func main() {

	var n int
	fmt.Scan(&n)

	a, b := 1, 1

	for i := 1; i <= n; i++ {

		fmt.Println(a)

		t := a + b
		a = b
		b = t

	}

}
