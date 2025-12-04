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
	input = strings.TrimSpace(input)
	grid := make([][]rune, len(strings.Split(input, "\n")))
	for i := range grid {
		grid[i] = []rune(strings.Split(input, "\n")[i])
	}

	rows, cols := len(grid), len(grid[0])
	result := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '@' {
				if adjacentCountValid(grid, i, j) {
					result += 1
				}
			}
		}
	}
	return result
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
