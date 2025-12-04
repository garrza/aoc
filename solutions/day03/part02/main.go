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
		bankJoltage := calculateBankJoltage(bank)
		total += bankJoltage
	}
	return total
}

func calculateBankJoltage(bank string) int {
	// we need to select exactly 12 digits to form the largest possible number.
	// since we want the largest number, we want the largest digits as early as possible.

	var resultBuilder strings.Builder
	currentIdx := 0
	digitsNeeded := 12

	for i := 0; i < 12; i++ {
		remainingNeeded := digitsNeeded - 1 - i
		// the furthest we can go and still have enough digits left
		searchLimit := len(bank) - remainingNeeded

		bestDigit := -1
		bestDigitIdx := -1

		// search for the largest digit in the valid window
		for j := currentIdx; j < searchLimit; j++ {
			digit := int(bank[j] - '0')
			if digit == 9 { // optimization: 9 is the max possible, take it immediately
				bestDigit = 9
				bestDigitIdx = j
				break
			}
			if digit > bestDigit {
				bestDigit = digit
				bestDigitIdx = j
			}
		}

		// append the best digit found
		resultBuilder.WriteString(strconv.Itoa(bestDigit))

		// advance our start pointer to just after the digit we picked
		currentIdx = bestDigitIdx + 1
	}

	// convert the 12-digit string to an integer
	val, err := strconv.Atoi(resultBuilder.String())
	if err != nil {
		log.Fatalf("failed to convert result to int: %v", err)
	}
	return val
}
