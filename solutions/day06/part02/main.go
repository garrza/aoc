package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	// Parse columns from right to left, where each column represents a digit position
	// Blank columns (no digits) separate different problems
	// Example:
	// 	123 328  51 64
	//  45 64  387 23
	//   6 98  215 314
	//  *   +   *   +
	// Rightmost problem: columns [4,3,2] -> 4 + 431 + 623 = 1058
	// Next problem: columns [5,1,7] -> 175 * 581 * 32 = 3253600

	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	if len(lines) == 0 {
		return 0
	}

	operators := strings.Fields(lines[len(lines)-1])
	matrixLines := lines[:len(lines)-1]

	// find max width across all lines
	maxWidth := 0
	for _, line := range matrixLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// prse columns right-to-left, collecting problems
	var problems [][]int
	var currentProblem []int

	for col := maxWidth - 1; col >= 0; col-- {
		columnDigits := extractColumn(matrixLines, col)

		if len(columnDigits) > 0 {
			// Column has digits - add number to current problem
			num, _ := strconv.Atoi(columnDigits)
			currentProblem = append(currentProblem, num)
		} else if len(currentProblem) > 0 {
			// Blank column acts as separator - save current problem
			problems = append(problems, currentProblem)
			currentProblem = []int{}
		}
	}

	// last problem
	if len(currentProblem) > 0 {
		problems = append(problems, currentProblem)
	}

	// evaluate each problem with its operator (problems are in reverse order)
	result := 0
	for i, problem := range problems {
		op := ""
		if i < len(operators) {
			op = operators[len(operators)-1-i]
		}
		result += evaluateExpression(problem, op)
	}

	return result
}

// extractColumn reads a column top-to-bottom and returns the digits as a string
func extractColumn(lines []string, col int) string {
	digits := ""
	for _, line := range lines {
		if col < len(line) {
			char := line[col]
			if char >= '0' && char <= '9' {
				digits += string(char)
			}
		}
	}
	return digits
}

func evaluateExpression(col []int, operator string) int {
	switch operator {
	case "+":
		return sum(col)
	case "*":
		return product(col)
	}
	return 0
}

func sum(col []int) int {
	sum := 0
	for _, num := range col {
		sum += num
	}
	return sum
}

func product(col []int) int {
	product := 1
	for _, num := range col {
		product *= num
	}
	return product
}
