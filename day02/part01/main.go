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
	// id ranges into two parts: the start and the end
	// these ranges are separated by '-'
	// all of the id ranges are separated by ','
	sum := 0
	for _, idRange := range strings.Split(input, ",") {
		// int can't have leading zeros, converting to int makes sure we don't have any invalid invalid ids
		start, end, _ := strings.Cut(idRange, "-")
		startInt, err := strconv.Atoi(start)
		if err != nil {
			panic(err)
		}
		endInt, err := strconv.Atoi(end)
		if err != nil {
			panic(err)
		}
		// we now want to check if the id is valid
		for i := startInt; i <= endInt; i++ {
			if isInvalid(i) {
				sum += i
			}
		}
	}
	return sum
}

func isInvalid(id int) bool {
	// an id is invalid when it has a sequence of digits repeated twice (55 (5 twice), 6464 (64 twice), 123123 (123 twice))
	// none of the numbers have leading zeros (0101 isn't an ID at all, 101 is a valid ID that you would ignore)
	// 38593859 would be invalid because 38593859 has a sequence of digits repeated twice (3859)
	idStr := strconv.Itoa(id)
	if len(idStr)%2 != 0 {
		return false
	}
	half := len(idStr) / 2
	return idStr[:half] == idStr[half:]
}
