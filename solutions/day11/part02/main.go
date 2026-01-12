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

	// count all paths from "svr" to "out" that visit both "dac" and "fft"
	memo := make(map[string]int)
	pathCount := countPathsWithRequired(graph, "svr", "out", make(map[string]bool), false, false, memo)

	fmt.Println(pathCount)
}

// state key for memoization
func stateKey(node string, foundDac, foundFft bool) string {
	return fmt.Sprintf("%s_%t_%t", node, foundDac, foundFft)
}

// countPathsWithRequired uses dfs to count paths from start to end that visit both dac and fft
func countPathsWithRequired(graph map[string][]string, start, end string, visited map[string]bool, foundDac, foundFft bool, memo map[string]int) int {
	// check if current node is dac or fft
	if start == "dac" {
		foundDac = true
	}
	if start == "fft" {
		foundFft = true
	}

	// if we reached the destination, check if we visited both required nodes
	if start == end {
		if foundDac && foundFft {
			return 1
		}
		return 0
	}

	// check memoization (only if not in a cycle)
	key := stateKey(start, foundDac, foundFft)
	if !visited[start] {
		if val, ok := memo[key]; ok {
			return val
		}
	}

	// mark current node as visited
	visited[start] = true
	defer delete(visited, start) // unmark when backtracking

	totalPaths := 0

	// explore all neighbors
	for _, neighbor := range graph[start] {
		// only visit if not already in current path (avoid cycles)
		if !visited[neighbor] {
			totalPaths += countPathsWithRequired(graph, neighbor, end, visited, foundDac, foundFft, memo)
		}
	}

	// memoize result
	memo[key] = totalPaths

	return totalPaths
}
