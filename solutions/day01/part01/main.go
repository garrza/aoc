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
		distanceInt, err := strconv.Atoi(instruction[1:])
		if err != nil {
			panic(err)
		}
		if instruction[0] == 'L' {
			pos = (pos - distanceInt) % 100
			if pos == 0 {
				zeroCount++
			}
		} else {
			pos = (pos + distanceInt) % 100
			if pos == 0 {
				zeroCount++
			}
		}
	}
	return zeroCount
}
