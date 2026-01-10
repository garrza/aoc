package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	numLights int
	target    []bool
	buttons   [][]int
}

func parseLine(line string) Machine {
	// parse indicator lights [.##.]
	startBracket := strings.Index(line, "[")
	endBracket := strings.Index(line, "]")
	lightPattern := line[startBracket+1 : endBracket]

	numLights := len(lightPattern)
	target := make([]bool, numLights)
	for i, ch := range lightPattern {
		target[i] = (ch == '#')
	}

	// parse buttons (1,3,5) (0,2)...
	buttons := [][]int{}
	rest := line[endBracket+1:]

	for {
		startParen := strings.Index(rest, "(")
		if startParen == -1 {
			break
		}
		endParen := strings.Index(rest[startParen:], ")")
		if endParen == -1 {
			break
		}

		buttonStr := rest[startParen+1 : startParen+endParen]
		parts := strings.Split(buttonStr, ",")
		button := []int{}
		for _, p := range parts {
			num, _ := strconv.Atoi(strings.TrimSpace(p))
			button = append(button, num)
		}
		buttons = append(buttons, button)
		rest = rest[startParen+endParen+1:]
	}

	return Machine{
		numLights: numLights,
		target:    target,
		buttons:   buttons,
	}
}

// gaussian elimination over gf(2) to solve the system and find minimum solution
func solveGF2(machine Machine) int {
	n := machine.numLights
	m := len(machine.buttons)

	// create augmented matrix [a | b]
	// a[i][j] = 1 if button j toggles light i
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m+1)
		if machine.target[i] {
			matrix[i][m] = 1 // target value
		}
	}

	// fill in which buttons toggle which lights
	for j, button := range machine.buttons {
		for _, light := range button {
			matrix[light][j] = 1
		}
	}

	// gaussian elimination to rref
	pivotCols := []int{}
	pivotRow := 0
	for col := 0; col < m && pivotRow < n; col++ {
		// find pivot
		foundPivot := false
		for row := pivotRow; row < n; row++ {
			if matrix[row][col] == 1 {
				// swap rows
				matrix[pivotRow], matrix[row] = matrix[row], matrix[pivotRow]
				foundPivot = true
				break
			}
		}

		if !foundPivot {
			continue
		}

		pivotCols = append(pivotCols, col)

		// eliminate
		for row := 0; row < n; row++ {
			if row != pivotRow && matrix[row][col] == 1 {
				for c := 0; c <= m; c++ {
					matrix[row][c] ^= matrix[pivotRow][c]
				}
			}
		}
		pivotRow++
	}

	// check for inconsistency
	for row := 0; row < n; row++ {
		allZero := true
		for col := 0; col < m; col++ {
			if matrix[row][col] == 1 {
				allZero = false
				break
			}
		}
		if allZero && matrix[row][m] == 1 {
			// inconsistent system - no solution
			return -1
		}
	}

	// find free variables (columns without pivots)
	isPivot := make([]bool, m)
	for _, col := range pivotCols {
		isPivot[col] = true
	}

	freeVars := []int{}
	for col := 0; col < m; col++ {
		if !isPivot[col] {
			freeVars = append(freeVars, col)
		}
	}

	// try all combinations of free variables to find minimum solution
	minCount := m + 1
	numFree := len(freeVars)

	for mask := 0; mask < (1 << numFree); mask++ {
		solution := make([]int, m)

		// set free variables according to mask
		for i, freeVar := range freeVars {
			if (mask & (1 << i)) != 0 {
				solution[freeVar] = 1
			}
		}

		// compute pivot variables
		for i := len(pivotCols) - 1; i >= 0; i-- {
			pivotCol := pivotCols[i]
			row := i

			val := matrix[row][m]
			for col := pivotCol + 1; col < m; col++ {
				if matrix[row][col] == 1 {
					val ^= solution[col]
				}
			}
			solution[pivotCol] = val
		}

		// count button presses
		count := 0
		for _, val := range solution {
			count += val
		}

		if count < minCount {
			minCount = count
		}
	}

	return minCount
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalPresses := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		machine := parseLine(line)
		presses := solveGF2(machine)
		if presses == -1 {
			fmt.Println("no solution found for machine")
			return
		}
		totalPresses += presses
	}

	fmt.Println(totalPresses)
}
