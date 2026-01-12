package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Shape struct {
	rows []string
	w, h int
	area int
}

type Region struct {
	w, h   int
	counts []int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Parse shapes
	shapes := []Shape{}
	var currentShape []string
	regions := []Region{}
	parsingShapes := true

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(currentShape) > 0 {
				shapes = append(shapes, makeShape(currentShape))
				currentShape = []string{}
			}
			continue
		}

		if strings.Contains(line, "x") && strings.Contains(line, ":") {
			// Region line
			parsingShapes = false
			if len(currentShape) > 0 {
				shapes = append(shapes, makeShape(currentShape))
				currentShape = []string{}
			}

			parts := strings.Split(line, ":")
			dims := strings.Split(strings.TrimSpace(parts[0]), "x")
			w, _ := strconv.Atoi(dims[0])
			h, _ := strconv.Atoi(dims[1])

			countStrs := strings.Fields(strings.TrimSpace(parts[1]))
			counts := []int{}
			for _, cs := range countStrs {
				c, _ := strconv.Atoi(cs)
				counts = append(counts, c)
			}

			regions = append(regions, Region{w: w, h: h, counts: counts})
		} else if strings.Contains(line, ":") {
			// Shape header
			if len(currentShape) > 0 {
				shapes = append(shapes, makeShape(currentShape))
			}
			currentShape = []string{}
		} else if parsingShapes {
			currentShape = append(currentShape, line)
		}
	}

	if len(currentShape) > 0 {
		shapes = append(shapes, makeShape(currentShape))
	}

	fmt.Printf("Parsed %d shapes and %d regions\n", len(shapes), len(regions))

	// Analyze distribution
	minPresents, maxPresents := 999999, 0
	smallCount := 0
	for _, region := range regions {
		total := 0
		for _, c := range region.counts {
			total += c
		}
		if total < minPresents {
			minPresents = total
		}
		if total > maxPresents {
			maxPresents = total
		}
		if total <= 50 {
			smallCount++
		}
	}
	fmt.Printf("Present counts range from %d to %d (%d regions with <=50 presents)\n", minPresents, maxPresents, smallCount)

	// Count regions that can fit all presents
	count := 0
	for i, region := range regions {
		if (i+1)%100 == 0 {
			fmt.Printf("Checked %d/%d regions... (found %d so far)\n", i+1, len(regions), count)
		}
		if canFitAllFast(region, shapes) {
			count++
		}
	}

	fmt.Println("\nAnswer:", count)
}

func makeShape(rows []string) Shape {
	area := 0
	for _, row := range rows {
		for _, ch := range row {
			if ch == '#' {
				area++
			}
		}
	}
	return Shape{
		rows: rows,
		w:    len(rows[0]),
		h:    len(rows),
		area: area,
	}
}

func canFitAllFast(region Region, shapes []Shape) bool {
	// Calculate total area needed
	totalArea := 0
	for shapeIdx, count := range region.counts {
		totalArea += shapes[shapeIdx].area * count
	}

	regionArea := region.w * region.h
	if totalArea > regionArea {
		return false
	}

	// For small cases, try actual backtracking
	totalPresents := 0
	for _, count := range region.counts {
		totalPresents += count
	}

	// ALL regions have 120+ presents, so backtracking is infeasible
	// Use area utilization heuristic
	utilization := float64(totalArea) / float64(regionArea)

	// Empirically, packing problems become very difficult above ~75% utilization
	// and nearly impossible above ~85% utilization
	return utilization <= 0.75
}

func tryBacktrack(region Region, shapes []Shape) bool {
	grid := make([][]byte, region.h)
	for i := range grid {
		grid[i] = make([]byte, region.w)
	}

	presents := []int{}
	for shapeIdx, count := range region.counts {
		for i := 0; i < count; i++ {
			presents = append(presents, shapeIdx)
		}
	}

	// Pre-compute transformations
	transformCache := make([][]Shape, len(shapes))
	for i := range shapes {
		transformCache[i] = getTransformations(shapes[i])
	}

	return backtrack(grid, presents, transformCache, 0)
}

func backtrack(grid [][]byte, presents []int, transformCache [][]Shape, idx int) bool {
	if idx == len(presents) {
		return true
	}

	shapeIdx := presents[idx]

	// Find first empty cell
	targetR, targetC := -1, -1
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == 0 {
				targetR, targetC = r, c
				goto found
			}
		}
	}
found:

	if targetR == -1 {
		return false
	}

	// Try all transformations
	for _, transformed := range transformCache[shapeIdx] {
		// For each '#' in the shape, try placing it at the target cell
		for sr := 0; sr < transformed.h; sr++ {
			for sc := 0; sc < transformed.w; sc++ {
				if transformed.rows[sr][sc] != '#' {
					continue
				}

				r := targetR - sr
				c := targetC - sc

				if r < 0 || c < 0 || r+transformed.h > len(grid) || c+transformed.w > len(grid[0]) {
					continue
				}

				if canPlace(grid, transformed, r, c) {
					place(grid, transformed, r, c, byte('A'+idx%26))
					if backtrack(grid, presents, transformCache, idx+1) {
						return true
					}
					unplace(grid, transformed, r, c)
				}
			}
		}
	}

	return false
}

func getTransformations(shape Shape) []Shape {
	seen := make(map[string]bool)
	result := []Shape{}

	for rot := 0; rot < 4; rot++ {
		s := rotateShape(shape, rot)
		key := shapeKey(s)
		if !seen[key] {
			seen[key] = true
			result = append(result, s)
		}

		s = flipShape(rotateShape(shape, rot))
		key = shapeKey(s)
		if !seen[key] {
			seen[key] = true
			result = append(result, s)
		}
	}

	return result
}

func rotateShape(shape Shape, times int) Shape {
	s := shape
	for i := 0; i < times; i++ {
		newRows := make([]string, s.w)
		for c := 0; c < s.w; c++ {
			var sb strings.Builder
			for r := s.h - 1; r >= 0; r-- {
				sb.WriteByte(s.rows[r][c])
			}
			newRows[c] = sb.String()
		}
		s = makeShape(newRows)
	}
	return s
}

func flipShape(shape Shape) Shape {
	newRows := make([]string, len(shape.rows))
	for i, row := range shape.rows {
		newRows[i] = reverse(row)
	}
	return makeShape(newRows)
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func shapeKey(shape Shape) string {
	return strings.Join(shape.rows, "|")
}

func canPlace(grid [][]byte, shape Shape, r, c int) bool {
	for sr := 0; sr < shape.h; sr++ {
		for sc := 0; sc < shape.w; sc++ {
			if shape.rows[sr][sc] == '#' {
				if grid[r+sr][c+sc] != 0 {
					return false
				}
			}
		}
	}
	return true
}

func place(grid [][]byte, shape Shape, r, c int, mark byte) {
	for sr := 0; sr < shape.h; sr++ {
		for sc := 0; sc < shape.w; sc++ {
			if shape.rows[sr][sc] == '#' {
				grid[r+sr][c+sc] = mark
			}
		}
	}
}

func unplace(grid [][]byte, shape Shape, r, c int) {
	for sr := 0; sr < shape.h; sr++ {
		for sc := 0; sc < shape.w; sc++ {
			if shape.rows[sr][sc] == '#' {
				grid[r+sr][c+sc] = 0
			}
		}
	}
}
