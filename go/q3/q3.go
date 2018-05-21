package main

import (
	"fmt"
	"os"
	"strconv"
)

func q3(n int) {
	num := (n*n + 1) / 2
	fmt.Printf("[ %d", num)

	// each value of i represents moving in one direction (right,down,left,up) as far as needed.
	for i := 1; i < 2*n; i++ {
		steps := (i + 1) / 2 // steps in current direction -> [1,1,2,2,3,3,...]

		for j := 1; j <= steps && j < n; j++ {
			num += [4]int{-n, 1, n, -1}[i%4] // [up,right,down,left] -> num + [-N,1,N,-1]
			fmt.Printf(", %d", num)
		}
	}
	fmt.Printf(", end]")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: q3 N")
		fmt.Println("where N is the width/length of the grid of numbers.")
		os.Exit(1)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if n%2 == 0 {
		fmt.Println("N must be an odd number.")
		os.Exit(1)
	}

	fmt.Printf("Running for %dx%d matrix.\n", n, n)
	q3(n)
}
