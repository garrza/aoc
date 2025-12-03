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

// isRepeatedSequence checks if the ID consists of a sequence of digits repeated twice or more.
func isInvalidId(id int) bool {
	s := strconv.Itoa(id)
	n := len(s)

	// Try all possible lengths for the repeating sequence.
	// The sequence length `k` must be at most n/2 because it needs to repeat at least twice.
	for k := 1; k <= n/2; k++ {
		// If the total length isn't divisible by k, this sequence length can't work.
		if n%k != 0 {
			continue
		}

		pattern := s[:k]
		// Check if repeating this pattern fills the string.
		if strings.Repeat(pattern, n/k) == s {
			return true
		}
	}

	return false
}
