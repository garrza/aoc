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
	// batteries labed with their joltage rating (1-9)
	// each line represents a 'bank', we want to combine two digits such that we get the highest possible number
	// after getting each bank, we sum the digits and add to the total
	// we return the total
	input = strings.TrimSpace(input)

	total := 0
	for _, bank := range strings.Split(input, "\n") {
		bankJoltage := calculateBankJoltageOptimized(bank)
		total += bankJoltage
	}
	return total
}

func calculateBankJoltage(bank string) int {
	// we need to find two digits in the string (at indices i and j where i < j)
	// such that the number formed by concatenating them is maximized.
	// since we want to maximize a 2-digit number, we prioritize the first digit (tens place).

	maxJoltage := 0

	// iterate through all possible pairs
	for i := 0; i < len(bank); i++ {
		for j := i + 1; j < len(bank); j++ {
			digit1 := int(bank[i] - '0')
			digit2 := int(bank[j] - '0')

			joltage := digit1*10 + digit2
			if joltage > maxJoltage {
				maxJoltage = joltage
			}
		}
	}

	return maxJoltage
}

func calculateBankJoltageOptimized(bank string) int {
	// we need to find two digits in the string (at indices i and j where i < j)
	// such that the number formed by concatenating them is maximized.
	// since we want to maximize a 2-digit number, we prioritize the first digit (tens place).
	// we can reduce complexity to o(n) time by precomputing the maximum digit to the right of each index

	n := len(bank)
	if n < 2 {
		return 0
	}

	maxSuffix := make([]int, n)
	currentMax := -1

	for i := n - 1; i >= 0; i-- {
		digit := int(bank[i] - '0')
		if digit > currentMax {
			currentMax = digit
		}
		maxSuffix[i] = currentMax
	}

	maxJoltage := 0

	// the second digit is the maximum digit to the right of the current index
	for i := 0; i < n-1; i++ {
		digit1 := int(bank[i] - '0')
		digit2 := maxSuffix[i+1]
		joltage := digit1*10 + digit2
		if joltage > maxJoltage {
			maxJoltage = joltage
		}
	}

	return maxJoltage
}
