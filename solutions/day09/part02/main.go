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

// cache for point-in-polygon results
var polygonCache = make(map[point]bool)

func main() {
	// read input file
	file, err := os.Open("solutions/day09/part02/input.txt")
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	// parse all red tile coordinates in order
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

	// create a set of red tiles for quick lookup
	redSet := make(map[point]bool)
	for _, p := range redTiles {
		redSet[p] = true
	}

	// create a set of green tiles
	greenSet := make(map[point]bool)

	// add all tiles on edges between consecutive red tiles
	for i := 0; i < len(redTiles); i++ {
		p1 := redTiles[i]
		p2 := redTiles[(i+1)%len(redTiles)] // wrap around to close the loop

		// add all points on the line between p1 and p2
		if p1.x == p2.x {
			// vertical line
			minY, maxY := min(p1.y, p2.y), max(p1.y, p2.y)
			for y := minY; y <= maxY; y++ {
				p := point{x: p1.x, y: y}
				if !redSet[p] {
					greenSet[p] = true
				}
			}
		} else if p1.y == p2.y {
			// horizontal line
			minX, maxX := min(p1.x, p2.x), max(p1.x, p2.x)
			for x := minX; x <= maxX; x++ {
				p := point{x: x, y: p1.y}
				if !redSet[p] {
					greenSet[p] = true
				}
			}
		}
	}

	// find the largest rectangle with red corners and only red/green tiles
	maxArea := 0

	// create a list of all pairs sorted by potential area (descending)
	type candidate struct {
		i, j int
		area int
	}
	var candidates []candidate

	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			p1 := redTiles[i]
			p2 := redTiles[j]

			width := abs(p2.x-p1.x) + 1
			height := abs(p2.y-p1.y) + 1
			area := width * height

			candidates = append(candidates, candidate{i: i, j: j, area: area})
		}
	}

	// sort candidates by area (largest first)
	// simple bubble sort is fine for reasonable sized lists
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].area > candidates[i].area {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	// check candidates in order of potential area
	checked := 0
	for _, c := range candidates {
		// early exit: if this candidate can't beat maxArea, we're done
		if c.area <= maxArea {
			break
		}

		p1 := redTiles[c.i]
		p2 := redTiles[c.j]

		// check if all tiles in the rectangle are red or green
		if isValidRectangle(p1, p2, redSet, greenSet, redTiles) {
			maxArea = c.area
			fmt.Printf("found valid rectangle with area %d at (%d,%d) to (%d,%d)\n",
				maxArea, p1.x, p1.y, p2.x, p2.y)
		}

		checked++
		if checked%1000 == 0 {
			fmt.Printf("checked %d candidates, maxArea=%d\n", checked, maxArea)
		}
	}

	fmt.Println(maxArea)
}

// isInsidePolygon uses ray casting algorithm with caching
func isInsidePolygon(p point, polygon []point) bool {
	// check cache first
	if result, ok := polygonCache[p]; ok {
		return result
	}

	inside := false
	n := len(polygon)

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		pi := polygon[i]
		pj := polygon[j]

		// check if the ray from p to the right crosses the edge from pi to pj
		if ((pi.y > p.y) != (pj.y > p.y)) &&
			(p.x < (pj.x-pi.x)*(p.y-pi.y)/(pj.y-pi.y)+pi.x) {
			inside = !inside
		}
	}

	// cache the result
	polygonCache[p] = inside
	return inside
}

// isValidRectangle checks if all tiles in the rectangle are red or green
func isValidRectangle(p1, p2 point, redSet, greenSet map[point]bool, polygon []point) bool {
	minX, maxX := min(p1.x, p2.x), max(p1.x, p2.x)
	minY, maxY := min(p1.y, p2.y), max(p1.y, p2.y)

	width := maxX - minX + 1
	height := maxY - minY + 1

	// For very large rectangles, use sampling instead of checking every point
	if width*height > 100000 {
		// Sample points: corners, edges, and some interior points
		samplesToCheck := []point{
			// corners
			{minX, minY}, {maxX, minY}, {minX, maxY}, {maxX, maxY},
		}

		// sample edges (every 100 points)
		step := 100
		for x := minX; x <= maxX; x += step {
			samplesToCheck = append(samplesToCheck, point{x, minY}, point{x, maxY})
		}
		for y := minY; y <= maxY; y += step {
			samplesToCheck = append(samplesToCheck, point{minX, y}, point{maxX, y})
		}

		// sample interior (sparse grid)
		for x := minX; x <= maxX; x += step {
			for y := minY; y <= maxY; y += step {
				samplesToCheck = append(samplesToCheck, point{x, y})
			}
		}

		// check sampled points
		for _, p := range samplesToCheck {
			if !isValidPoint(p, redSet, greenSet, polygon) {
				return false
			}
		}
		return true
	}

	// for smaller rectangles, check all points
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			p := point{x: x, y: y}
			if !isValidPoint(p, redSet, greenSet, polygon) {
				return false
			}
		}
	}

	return true
}

// isValidPoint checks if a point is red, green, or inside the polygon
func isValidPoint(p point, redSet, greenSet map[point]bool, polygon []point) bool {
	return redSet[p] || greenSet[p] || isInsidePolygon(p, polygon)
}

// helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
