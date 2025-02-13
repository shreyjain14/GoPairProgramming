package main

import "fmt"

func main() {
	var a int
	var sum int
	fmt.Scanf("%d", &a)
	for i := 0; i < a; i++ {
		sum += i
	}
	fmt.Println(sum)
}
