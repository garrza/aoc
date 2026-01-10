package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// point represents a coordinate on the grid
type point struct {
	x, y int
}

func main() {
	// read input file
	file, err := os.Open("solutions/day09/part01/input.txt")
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	// parse all red tile coordinates
	var redTiles []point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// parse x,y format
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		y, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		redTiles = append(redTiles, point{x: x, y: y})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
		return
	}

	// find the largest rectangle
	maxArea := 0

	// try every pair of red tiles as opposite corners
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			p1 := redTiles[i]
			p2 := redTiles[j]

			// calculate rectangle area using these two points as opposite corners
			// add 1 to each dimension to count tiles inclusively
			width := abs(p2.x-p1.x) + 1
			height := abs(p2.y-p1.y) + 1
			area := width * height

			if area > maxArea {
				maxArea = area
			}
		}
	}

	fmt.Println(maxArea)
}

// abs returns the absolute value of an integer
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
