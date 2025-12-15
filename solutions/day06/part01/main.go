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
	// we will be given input as rows of numbers, separated by spaces
	// the final row includes an operator, either + or *
	// we need to evaluate the expression given the operator for all columns
	// we must then add the results of all of our operations, that is the output
	// the idea will be to create a matrix of numbers, and then evaluate the expression for each column

	lines := strings.Split(strings.TrimSpace(input), "\n")
	operators := strings.Fields(lines[len(lines)-1])

	matrixLines := lines[:len(lines)-1]
	matrix := make([][]int, len(matrixLines))
	for i, line := range matrixLines {
		fields := strings.Fields(line)
		matrix[i] = make([]int, len(fields))
		for j, num := range fields {
			matrix[i][j], _ = strconv.Atoi(num)
		}
	}

	result := 0
	if len(matrix) > 0 {
		for i := 0; i < len(matrix[0]); i++ {
			col := make([]int, len(matrix))
			for j := 0; j < len(matrix); j++ {
				col[j] = matrix[j][i]
			}

			op := ""
			if i < len(operators) {
				op = operators[i]
			}
			result += evaluateExpression(col, op)
		}
	}
	return result

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
