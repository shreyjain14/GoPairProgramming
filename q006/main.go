package main

import "fmt"

func main() {

	var n int
	fmt.Scanf("%d", &n)

	x := n
	if x == 0 {
		fmt.Println("0")
	} else {
		result := ""
		for x > 0 {
			if x%2 == 0 {
				result = "0" + result
			} else {
				result = "1" + result
			}
			x /= 2
		}
		fmt.Println(result)
	}
}
