package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	// build the graph (adjacency list)
	graph := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// parse line: "device: output1 output2 ..."
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		device := strings.TrimSpace(parts[0])
		outputsStr := strings.TrimSpace(parts[1])
		outputs := strings.Fields(outputsStr)

		graph[device] = outputs
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
		return
	}

	// count all paths from "you" to "out"
	pathCount := countPaths(graph, "you", "out", make(map[string]bool))

	fmt.Println(pathCount)
}

// countpaths uses dfs to count all paths from start to end
func countPaths(graph map[string][]string, start, end string, visited map[string]bool) int {
	// if we reached the destination, we found one path
	if start == end {
		return 1
	}

	// mark current node as visited
	visited[start] = true
	defer delete(visited, start) // unmark when backtracking

	totalPaths := 0

	// explore all neighbors
	for _, neighbor := range graph[start] {
		// only visit if not already in current path (avoid cycles)
		if !visited[neighbor] {
			totalPaths += countPaths(graph, neighbor, end, visited)
		}
	}

	return totalPaths
}
