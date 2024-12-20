package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Example")
	exampleLocation := "./input/example.txt"
	available, liked := readInputFile(exampleLocation)
	one, two := solve(available, liked)
	fmt.Println("Part One: ", one)
	fmt.Println("Part Two: ", two)

	fmt.Println("\nInput")
	inputLocation := "./input/input.txt"
	available, liked = readInputFile(inputLocation)
	one, two = solve(available, liked)
	fmt.Println("Part One: ", one)
	fmt.Println("Part Two: ", two)
}

func solve(available AvailablePatterns, liked LikedPatterns) (int, int) {
	var p1, p2 int
	var memo = make(map[string]int)
	for pattern := range liked {
		patt := string(pattern)

		ways := checkPattern(patt, available, memo)
		if ways > 0 {
			p1++
			p2 += ways
		}
	}

	return p1, p2
}

func checkPattern(pattern string, available AvailablePatterns, memo map[string]int) int {
	if pattern == "" {
		return 1
	}

	if val, exists := memo[pattern]; exists {
		return val
	}

	total := 0
	for availablePattern := range available {
		patStr := string(availablePattern)
		if strings.HasPrefix(pattern, patStr) {
			ways := checkPattern(pattern[len(patStr):], available, memo)
			total += ways
		}
	}

	memo[pattern] = total
	return total
}

type Pattern string
type AvailablePatterns map[Pattern]bool
type LikedPatterns map[Pattern]bool

func readInputFile(location string) (AvailablePatterns, LikedPatterns) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	availablePatterns := make(AvailablePatterns)

	scanner.Scan()
	line := scanner.Text()
	patterns := strings.FieldsFunc(line, func(r rune) bool {
		if r == ',' {
			return true
		}

		return false
	})

	for _, pattern := range patterns {
		availablePatterns[Pattern(strings.Trim(pattern, " "))] = true
	}

	likedPatterns := make(LikedPatterns)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		likedPatterns[Pattern(line)] = true
	}

	return availablePatterns, likedPatterns
}
