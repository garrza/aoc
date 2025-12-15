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
	// 	same thing, however, now we must also separate the digits themselves
	// 	this will be done from top to bottom, right-to-left
	// 	123 328  51 64
	//  45 64  387 23
	//   6 98  215 314
	//  *   +   *   +
	// the rightmost problem is 4 + 431 + 623 = 1058
	// the second problem from the right is 175 * 581 * 32 = 3253600

	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	if len(lines) == 0 {
		return 0
	}
	operatorLine := lines[len(lines)-1]
	operators := strings.Fields(operatorLine)

	matrixLines := lines[:len(lines)-1]

	// Find max width
	maxWidth := 0
	for _, line := range matrixLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	result := 0
	var currentBlock []int
	opIndex := len(operators) - 1

	for col := maxWidth - 1; col >= 0; col-- {
		// Extract digits for this column from top to bottom
		digits := ""
		for _, line := range matrixLines {
			if col < len(line) {
				char := line[col]
				if char >= '0' && char <= '9' {
					digits += string(char)
				}
			}
		}

		if len(digits) > 0 {
			num, _ := strconv.Atoi(digits)
			currentBlock = append(currentBlock, num)
		} else {
			// Empty column (separator)
			if len(currentBlock) > 0 {
				// Process the block
				op := ""
				if opIndex >= 0 {
					op = operators[opIndex]
					opIndex--
				}
				result += evaluateExpression(currentBlock, op)
				currentBlock = []int{}
			}
		}
	}

	// Process remaining block if any
	if len(currentBlock) > 0 {
		op := ""
		if opIndex >= 0 {
			op = operators[opIndex]
		}
		result += evaluateExpression(currentBlock, op)
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
