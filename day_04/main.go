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
	fmt.Println("PartTwo: ", partTwo(grid))

	fmt.Println("\nInput")
	location = "./input/input.txt"

	grid = readInputFile(location)

	fmt.Println("PartOne: ", partOne(grid))
	fmt.Println("PartTwo: ", partTwo(grid))
}

const WORD_P1 string = "XMAS"
const WORD_P2 string = "MAS"

func partOne(grid [][]byte) int {
	var res int = 0
	var startLetter byte = WORD_P1[0]
	var endLetter byte = WORD_P1[len(WORD_P1)-1]

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			current := grid[i][j]

			if current == startLetter || current == endLetter {
				res += checkWord(grid, i, j)
			}
		}
	}

	return res
}

func partTwo(grid [][]byte) int {
	var res int = 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			current := grid[i][j]

			if current == WORD_P2[1] {
				if checkWordP2(grid, i, j) {
					res++
				}
			}
		}
	}

	return res
}

func checkWordP2(grid [][]byte, i int, j int) bool {
	var rightX = i - 1
	var leftX = i + 1
	var topY = j - 1
	var bottomY = j + 1

	if rightX < 0 || leftX >= len(grid) || topY < 0 || bottomY >= len(grid[i]) {
		return false
	}

	// This is ugly, but it works. Couldn't be bothered to rotate grid or do fancy fors
	if (grid[rightX][topY] == WORD_P2[0] && grid[leftX][bottomY] == WORD_P2[2]) || (grid[rightX][topY] == WORD_P2[2] && grid[leftX][bottomY] == WORD_P2[0]) {
		if (grid[leftX][topY] == WORD_P2[0] && grid[rightX][bottomY] == WORD_P2[2]) || (grid[leftX][topY] == WORD_P2[2] && grid[rightX][bottomY] == WORD_P2[0]) {
			return true
		}
	}

	return false
}

func checkWord(grid [][]byte, i int, j int) int {
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
	for k := 0; k < len(WORD_P1); k++ {
		if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[i]) {
			return false
		}

		if grid[i][j] != WORD_P1[k] {
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
