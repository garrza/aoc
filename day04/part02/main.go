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

	inputStr := strings.TrimSpace(string(input))
	grid := make([][]rune, len(strings.Split(inputStr, "\n")))
	for i := range grid {
		grid[i] = []rune(strings.Split(inputStr, "\n")[i])
	}

	result := solve(grid, 0)
	if err != nil {
		log.Fatalf("failed to solve: %v", err)
	}

	fmt.Println(result)
}

func solve(grid [][]rune, totalRemovals int) int {
	rows, cols := len(grid), len(grid[0])
	var toRemove [][2]int

	// identify all removals for this step
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '@' {
				if adjacentCountValid(grid, i, j) {
					toRemove = append(toRemove, [2]int{i, j})
				}
			}
		}
	}

	// base case: no more removals possible
	if len(toRemove) == 0 {
		return totalRemovals
	}

	// apply removals
	for _, p := range toRemove {
		grid[p[0]][p[1]] = 'x'
	}

	// recursive step
	return solve(grid, totalRemovals+len(toRemove))
}

func adjacentCountValid(grid [][]rune, i, j int) bool {
	neighbours := 0
	rows, cols := len(grid), len(grid[0])
	dirs := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	for _, dir := range dirs {
		x, y := i+dir[0], j+dir[1]
		if x >= 0 && x < rows && y >= 0 && y < cols && grid[x][y] == '@' {
			neighbours++
		}
	}
	return neighbours < 4
}
