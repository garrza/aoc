package main

import (
	"fmt"
	"log"
	"os"
	"sort"
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

	// parse ranges
	type Range struct {
		start, end int
	}
	var ranges []Range
	// Use parts[0] which contains the ranges. If there's no double newline, parts[0] is the whole input.
	lines := strings.Split(parts[0], "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		rangeStr := strings.Split(line, "-")
		if len(rangeStr) != 2 {
			continue
		}
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

	// Sort ranges by start time
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	// Merge overlapping ranges
	var merged []Range
	if len(ranges) > 0 {
		merged = append(merged, ranges[0])
		for _, r := range ranges[1:] {
			last := &merged[len(merged)-1]
			// if current range overlaps with or touches the last range
			if r.start <= last.end+1 {
				if r.end > last.end {
					last.end = r.end
				}
			} else {
				merged = append(merged, r)
			}
		}
	}

	// Count unique ingredients
	freshCount := 0
	for _, r := range merged {
		freshCount += r.end - r.start + 1
	}

	return freshCount
}
