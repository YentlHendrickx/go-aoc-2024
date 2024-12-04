package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("\nExample")
	var location string = "./input/example.txt"

	grid := readInputFile(location)

	fmt.Println("PartOne: ", partOne(grid))
	// fmt.Println("PartTwo: ", solve(memSections, true))

	fmt.Println("\nInput")
	location = "./input/input.txt"

	grid = readInputFile(location)

	fmt.Println("PartOne: ", partOne(grid))
	// fmt.Println("PartTwo: ", solve(memSections, true))
}

const WORD string = "XMAS"

func partOne(grid [][]byte) int {
	var res int = 0
	var startLetter byte = WORD[0]
	var endLetter byte = WORD[len(WORD)-1]

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			current := grid[i][j]

			if current == startLetter || current == endLetter {
				// fmt.Printf("i: %d, j: %d, current: %c\n", i, j, current)
				res += checkWord(grid, i, j)
			}
		}
	}

	return res
}

func checkWord(grid [][]byte, i, j int) int {
	var count = 0

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if checkDirection(grid, i, j, x, y) {
				count++
			}
		}
	}

	return count
}

func checkDirection(grid [][]byte, i int, j int, x int, y int) bool {
	for k := 0; k < len(WORD); k++ {
		if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[i]) {
			return false
		}

		if grid[i][j] != WORD[k] {
			return false
		}
		i += x
		j += y
	}

	return true
}

func readInputFile(location string) [][]byte {
	file, err := os.Open(location)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result [][]byte

	for scanner.Scan() {
		result = append(result, []byte(scanner.Text()))
	}

	return result
}
