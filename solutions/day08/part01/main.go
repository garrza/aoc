package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Edge struct {
	i, j int
	dist float64
}

type UnionFind struct {
	parent []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		size:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // already connected
	}

	// union by size
	if uf.size[rootX] < uf.size[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	}

	return true
}

func (uf *UnionFind) GetComponentSizes() []int {
	components := make(map[int]int)
	for i := 0; i < len(uf.parent); i++ {
		root := uf.Find(i)
		components[root] = uf.size[root]
	}

	sizes := make([]int, 0, len(components))
	for _, size := range components {
		sizes = append(sizes, size)
	}
	return sizes
}

func distance(p1, p2 Point) float64 {
	dx := float64(p1.x - p2.x)
	dy := float64(p1.y - p2.y)
	dz := float64(p1.z - p2.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

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

	// parse points
	points := make([]Point, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		points[i] = Point{x, y, z}
	}

	// calculate all pairwise distances
	var edges []Edge
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := distance(points[i], points[j])
			edges = append(edges, Edge{i, j, dist})
		}
	}

	// sort edges by distance
	sort.Slice(edges, func(a, b int) bool {
		return edges[a].dist < edges[b].dist
	})

	// create union-find structure
	uf := NewUnionFind(len(points))

	// connect the 1000 closest pairs
	for i := 0; i < 1000 && i < len(edges); i++ {
		uf.Union(edges[i].i, edges[i].j)
	}

	// get all component sizes
	sizes := uf.GetComponentSizes()

	// sort sizes in descending order
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	// multiply the three largest
	if len(sizes) >= 3 {
		return sizes[0] * sizes[1] * sizes[2]
	}

	return 0
}
