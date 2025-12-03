package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}

	result, err := solve(string(input))
	if err != nil {
		log.Fatalf("failed to solve: %v", err)
	}

	fmt.Println(result)
}

func solve(input string) (int, error) {
	input = strings.TrimSpace(input)
	// Input format: "start-end,start-end,..."
	ranges := strings.Split(input, ",")

	sum := 0
	for _, r := range ranges {
		startStr, endStr, found := strings.Cut(r, "-")
		if !found {
			return 0, fmt.Errorf("invalid range format: %q", r)
		}

		start, err := strconv.Atoi(startStr)
		if err != nil {
			return 0, fmt.Errorf("invalid start value %q: %w", startStr, err)
		}

		end, err := strconv.Atoi(endStr)
		if err != nil {
			return 0, fmt.Errorf("invalid end value %q: %w", endStr, err)
		}

		for id := start; id <= end; id++ {
			if isInvalidId(id) {
				sum += id
			}
		}
	}

	return sum, nil
}

// isInvalidId checks if the ID consists of a sequence of digits repeated twice.
// Examples: 55 (5 repeated), 6464 (64 repeated), 123123 (123 repeated).
func isInvalidId(id int) bool {
	s := strconv.Itoa(id)
	n := len(s)

	// An ID must have an even number of digits to be split into two identical halves
	if n%2 != 0 {
		return false
	}

	half := n / 2
	return s[:half] == s[half:]
}
