package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	data := strings.TrimSpace(string(input))

	fmt.Println(solve(data))
}

func solve(input string) int {
	pos := 50

	zeroCount := 0
	for _, instruction := range strings.Split(input, "\n") {
		// we split the instruction into two parts: the direction and the distance
		direction := string(instruction[0])
		distance := string(instruction[1:])
		distanceInt, err := strconv.Atoi(distance)
		if err != nil {
			panic(err)
		}
		// we now want to count the number of times we pass through 0, new 0x434C49434B password method (clicks)
		for i := 0; i < distanceInt; i++ {
			if direction == "L" {
				pos = (pos - 1) % 100
			} else {
				pos = (pos + 1) % 100
			}
			if pos == 0 {
				zeroCount++
			}
		}
	}
	return zeroCount
}
