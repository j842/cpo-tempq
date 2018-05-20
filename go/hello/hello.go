package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func whee() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Welcome to the playground!")

	fmt.Println("The time is", time.Now())
	fmt.Println("My favorite number is", math.Sqrt(float64(rand.Intn(10))))
}

func main() {
	whee()
}
