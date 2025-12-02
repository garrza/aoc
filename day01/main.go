package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	data := strings.TrimSpace(string(input))

	fmt.Println("Part 1:", solvePart1(data))
	fmt.Println("Part 2:", solvePart2(data))
}

func solvePart1(input string) int {
	return 0
}

func solvePart2(input string) int {
	return 0
}

