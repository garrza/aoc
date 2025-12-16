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
	// given a grid
	// we must return how many times we 'split' our beam
	// beam only goes down, and is split by '^' characters
	// when a beam is split, it is split into two beams, one going left, one going right
	// we must return how many times we split the beam
	// beam always starts at the 'S' character

	// first, we need to create a grid from the input
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

	return calculateSplitCount(grid, startRow, startCol)
}

func calculateSplitCount(grid [][]rune, startRow, startCol int) int {
	type Beam struct {
		row, col int
	}

	q := []Beam{{startRow, startCol}}
	visit := make(map[string]bool)
	splitCount := 0

	for len(q) > 0 {
		beam := q[0]
		q = q[1:]

		// create a state to avoid reprocessing
		stateKey := fmt.Sprintf("%d,%d", beam.row, beam.col)
		if visit[stateKey] {
			continue
		}
		visit[stateKey] = true

		// all beams move down
		nextRow := beam.row + 1
		nextCol := beam.col

		// check if beam exits the grid
		if nextRow < 0 || nextRow >= len(grid) {
			continue
		}

		// check if the row is empty or column is out of bounds
		if len(grid[nextRow]) == 0 || nextCol < 0 || nextCol >= len(grid[nextRow]) {
			continue
		}

		cell := grid[nextRow][nextCol]
		if cell == '^' {
			splitCount++
			// create two new beams at immediate left and right of splitter
			// both beams continue downward from those positions
			leftCol := nextCol - 1
			rightCol := nextCol + 1

			if leftCol >= 0 {
				q = append(q, Beam{nextRow, leftCol})
			}
			if rightCol < len(grid[0]) {
				q = append(q, Beam{nextRow, rightCol})
			}
		} else {
			// beam continues downward
			q = append(q, Beam{nextRow, nextCol})
		}
	}

	return splitCount
}
