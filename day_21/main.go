package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Coordinate struct {
	x, y int
}

var numericPad = map[rune]Coordinate{
	'7': {0, 0}, '8': {1, 0}, '9': {2, 0},
	'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
	'1': {0, 2}, '2': {1, 2}, '3': {2, 2},
	'0': {1, 3}, 'A': {2, 3},
}

var directionalPad = map[rune]Coordinate{
	'^': {1, 0}, 'A': {2, 0},
	'<': {0, 1}, 'v': {1, 1}, '>': {2, 1},
}

func main() {
	fmt.Println("Example")
	solveBoth("./input/example.txt")

	fmt.Println("\nInput")
	solveBoth("./input/input.txt")
}

func solveBoth(location string) {
	fmt.Println("Part 1:", solve(location, 2))
	fmt.Println("Part 2:", solve(location, 25))
}

func solve(location string, layers int) int {
	codes := readInputFile(location)
	totalComplexity := 0
	cache := make(map[string][]int)

	for _, code := range codes {
		// Get initial sequence for numeric keypad
		baseSeq := getNumericPadSequence(code)

		// Calculate sequence length through all robot layers
		seqLen := getSequenceLength(baseSeq, layers, 1, cache)

		// Calculate numeric value and complexity
		numericValue := getNumericValue(code)
		complexity := seqLen * numericValue

		fmt.Printf("Code: %s, Sequence Length: %d, Numeric Value: %d, Complexity: %d\n",
			code, seqLen, numericValue, complexity)
		totalComplexity += complexity
	}

	return totalComplexity
}

func getNumericPadSequence(code string) []string {
	var sequence []string
	current := Coordinate{2, 3} // Start at 'A'

	for _, target := range code {
		dest := numericPad[target]

		// Calculate differences
		diffX := dest.x - current.x
		diffY := dest.y - current.y

		// Add horizontal movements
		for i := 0; i < abs(diffX); i++ {
			if diffX < 0 {
				sequence = append(sequence, "<")
			} else {
				sequence = append(sequence, ">")
			}
		}

		// Add vertical movements
		for i := 0; i < abs(diffY); i++ {
			if diffY < 0 {
				sequence = append(sequence, "v")
			} else {
				sequence = append(sequence, "^")
			}
		}

		// Add button press
		sequence = append(sequence, "A")
		current = dest
	}

	// Return to 'A'
	dest := numericPad['A']
	diffX := dest.x - current.x
	diffY := dest.y - current.y

	for i := 0; i < abs(diffX); i++ {
		if diffX < 0 {
			sequence = append(sequence, "<")
		} else {
			sequence = append(sequence, ">")
		}
	}

	for i := 0; i < abs(diffY); i++ {
		if diffY < 0 {
			sequence = append(sequence, "v")
		} else {
			sequence = append(sequence, "^")
		}
	}

	return sequence
}

func getDirectionalPadSequence(input []string) []string {
	var sequence []string
	current := Coordinate{2, 0} // Start at 'A'

	for _, move := range input {
		var target rune
		switch move {
		case "^":
			target = '^'
		case "v":
			target = 'v'
		case "<":
			target = '<'
		case ">":
			target = '>'
		case "A":
			target = 'A'
		}

		dest := directionalPad[target]
		diffX := dest.x - current.x
		diffY := dest.y - current.y

		// Prioritize horizontal movement for < and >
		if target == '<' || target == '>' {
			// Add horizontal movements
			for i := 0; i < abs(diffX); i++ {
				if diffX < 0 {
					sequence = append(sequence, "<")
				} else {
					sequence = append(sequence, ">")
				}
			}
			// Add vertical movements
			for i := 0; i < abs(diffY); i++ {
				if diffY < 0 {
					sequence = append(sequence, "v")
				} else {
					sequence = append(sequence, "^")
				}
			}
		} else {
			// For other buttons, prioritize vertical movement
			for i := 0; i < abs(diffY); i++ {
				if diffY < 0 {
					sequence = append(sequence, "v")
				} else {
					sequence = append(sequence, "^")
				}
			}
			for i := 0; i < abs(diffX); i++ {
				if diffX < 0 {
					sequence = append(sequence, "<")
				} else {
					sequence = append(sequence, ">")
				}
			}
		}

		sequence = append(sequence, "A")
		current = dest
	}

	return sequence
}

func getSequenceLength(input []string, maxLayers, currentLayer int, cache map[string][]int) int {
	key := strings.Join(input, "")

	// Check cache
	if val, ok := cache[key]; ok && val[currentLayer-1] != 0 {
		return val[currentLayer-1]
	}

	// Ensure cache entry exists
	if _, ok := cache[key]; !ok {
		cache[key] = make([]int, maxLayers)
	}

	// Get sequence for current layer
	seq := getDirectionalPadSequence(input)
	cache[key][0] = len(seq)

	if currentLayer == maxLayers {
		return len(seq)
	}

	// Split sequence into individual button presses
	var total int
	current := []string{}
	for _, move := range seq {
		current = append(current, move)
		if move == "A" {
			total += getSequenceLength(current, maxLayers, currentLayer+1, cache)
			current = []string{}
		}
	}

	cache[key][currentLayer-1] = total
	return total
}

func getNumericValue(code string) int {
	code = strings.TrimSuffix(code, "A")
	code = strings.TrimLeft(code, "0")

	if code == "" {
		return 0
	}

	value, err := strconv.Atoi(code)
	if err != nil {
		fmt.Printf("Error parsing numeric value from code '%s': %v\n", code, err)
		return 0
	}
	return value
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func readInputFile(location string) []string {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	var codes []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}

	return codes
}
