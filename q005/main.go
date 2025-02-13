package main

import "fmt"

func main() {

	var a string

	fmt.Scanf("%s", &a)

	l, r := 0, len(a)-1

	for l < r {

		if a[l] != a[r] {
			fmt.Println("NO")
			return
		}

		l++
		r--

	}

	fmt.Println("YES")

}
