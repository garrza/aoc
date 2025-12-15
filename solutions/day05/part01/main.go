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

	result := solve(string(input))
	fmt.Println(result)
}

func solve(input string) int {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	parts := strings.Split(input, "\n\n")
	if len(parts) < 2 {
		return 0
	}

	// parse ranges
	type Range struct {
		start, end int
	}
	var ranges []Range
	for _, line := range strings.Split(parts[0], "\n") {
		if line == "" {
			continue
		}
		rangeStr := strings.Split(line, "-")
		start, err := strconv.Atoi(rangeStr[0])
		if err != nil {
			log.Fatalf("failed to convert range start to int: %v", err)
		}
		end, err := strconv.Atoi(rangeStr[1])
		if err != nil {
			log.Fatalf("failed to convert range end to int: %v", err)
		}
		ranges = append(ranges, Range{start, end})
	}

	// count fresh ingredients
	freshCount := 0
	for _, line := range strings.Split(parts[1], "\n") {
		if line == "" {
			continue
		}
		ingredient, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("failed to convert ingredient to int: %v", err)
		}

		isFresh := false
		for _, r := range ranges {
			if ingredient >= r.start && ingredient <= r.end {
				isFresh = true
				break
			}
		}
		if isFresh {
			freshCount++
		}
	}

	return freshCount
}
