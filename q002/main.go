package main

import (
	"fmt"
	"math/rand"
)

func main() {
	randomNumber := rand.Intn(101)
	var guess int
	tries := 1
	fmt.Print("Guess the number: ")
	fmt.Scan(&guess)
	for guess != randomNumber {
		if guess > randomNumber {
			fmt.Println("Too high")
		} else {
			fmt.Println("Too low")
		}
		fmt.Print("Guess the number: ")
		fmt.Scan(&guess)
		tries++
	}
	fmt.Println("You Win. The random number was:", randomNumber, "and you took", tries, "tries.")
}
