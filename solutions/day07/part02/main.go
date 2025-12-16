package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}

	result := solve(string(input))
	if err != nil {
		log.Fatalf("failed to solve: %v", err)
	}

	fmt.Println(result)
}

func solve(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			grid = append(grid, []rune(line))
		}
	}

	var startRow, startCol int
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 'S' {
				startRow, startCol = i, j
				break
			}
		}
	}

	return countTimelines(grid, startRow, startCol)
}

func countTimelines(grid [][]rune, startRow, startCol int) int {
	// memoization: cache[row][col] = number of timelines from this position to exit
	cache := make(map[string]int)

	var dfs func(row, col int) int
	dfs = func(row, col int) int {
		key := fmt.Sprintf("%d,%d", row, col)
		if val, ok := cache[key]; ok {
			return val
		}

		// move down one row
		nextRow := row + 1

		// if we exit the grid, return 1 complete timeline
		if nextRow >= len(grid) {
			return 1 // exited the grid
		}

		if col < 0 || col >= len(grid[nextRow]) {
			return 1 // exited the grid
		}

		cell := grid[nextRow][col]
		var result int

		if cell == '^' {
			// splitter, particle takes BOTH paths (quantum superposition)
			// left path
			leftPaths := 0
			if col-1 >= 0 {
				leftPaths = dfs(nextRow, col-1)
			} else {
				leftPaths = 1 // exited left
			}

			// right path
			rightPaths := 0
			if col+1 < len(grid[nextRow]) {
				rightPaths = dfs(nextRow, col+1)
			} else {
				rightPaths = 1 // exited right
			}

			// total timelines = left timelines + right timelines
			result = leftPaths + rightPaths
		} else {
			// empty space: continue downward
			result = dfs(nextRow, col)
		}

		cache[key] = result
		return result
	}

	return dfs(startRow, startCol)
}
