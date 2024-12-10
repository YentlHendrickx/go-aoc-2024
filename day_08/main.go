package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Example")
	var location string = "./input/example.txt"
	grid := readInputFile(location)
	antennas := createAntennas(grid)

	one := solve(antennas, grid, false)
	fmt.Println("PartOne: ", one)

	two := solve(antennas, grid, true)
	fmt.Println("PartTwo: ", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	grid = readInputFile(location)
	antennas = createAntennas(grid)

	one = solve(antennas, grid, false)
	fmt.Println("PartOne: ", one)

	two = solve(antennas, grid, true)
	fmt.Println("PartTwo: ", two)
}

const PRINT_GRID = false

func solve(antennas Antennas, grid Grid, partTwo bool) int {
	var res int = 0
	var multLen int = 2

	if partTwo {
		multLen = len(grid)
	}

	var resonanceLocation map[int]map[int]bool = make(map[int]map[int]bool)
	for _, antennea := range antennas {
		for i := 0; i < len(antennea); i++ {
			for j := i + 1; j < len(antennea); j++ {
				x1, y1 := antennea[i][0], antennea[i][1]
				x2, y2 := antennea[j][0], antennea[j][1]

				// We can keep going as many times as needed
				for mult := 1; mult < multLen; mult++ {
					xDistance := (x1 - x2) * mult
					yDistance := (y1 - y2) * mult

					if xDistance > len(grid[0]) || yDistance > len(grid) {
						break
					}

					newX := x1 + xDistance
					newY := y1 + yDistance

					resonanceLocation = addResonance(newX, newY, resonanceLocation, grid)
					if partTwo {
						resonanceLocation = addResonance(x1, y1, resonanceLocation, grid)
					}

					newX = x2 - xDistance
					newY = y2 - yDistance

					resonanceLocation = addResonance(newX, newY, resonanceLocation, grid)
					if partTwo {
						resonanceLocation = addResonance(x2, y2, resonanceLocation, grid)
					}
				}
			}
		}
	}

	if PRINT_GRID {
		for y := range grid {
			for x := range grid[y] {
				if _, ok := resonanceLocation[x]; ok {
					if _, ok := resonanceLocation[x][y]; ok {
						fmt.Print("#")
					} else {
						fmt.Print(string(grid[y][x]))
					}
				} else {
					fmt.Print(string(grid[y][x]))
				}
			}
			fmt.Println()
		}
	}

	for i := range resonanceLocation {
		res += len(resonanceLocation[i])
	}

	return res
}

func addResonance(x, y int, resonanceLocation map[int]map[int]bool, grid Grid) map[int]map[int]bool {
	if !(x < len(grid[0]) && y < len(grid)) {
		return resonanceLocation
	} else if !(x >= 0 && y >= 0) {
		return resonanceLocation
	}

	if _, ok := resonanceLocation[x]; !ok {
		resonanceLocation[x] = make(map[int]bool)
	}

	if _, ok := resonanceLocation[y]; !ok {
		resonanceLocation[y] = make(map[int]bool)
	}

	resonanceLocation[x][y] = true

	return resonanceLocation
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func createAntennas(grid Grid) Antennas {
	var antennas Antennas = make(Antennas)

	for y, row := range grid {
		for x, cell := range row {
			if cell != '.' {
				if _, ok := antennas[cell]; !ok {
					antennas[cell] = make([][]int, 0)
				}

				antennas[cell] = append(antennas[cell], []int{x, y})
			}
		}
	}

	return antennas
}

type Antennas map[byte][][]int
type Grid [][]byte

func readInputFile(location string) Grid {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var out Grid = make(Grid, 0)

	for scanner.Scan() {
		line := scanner.Text()
		var row []byte = make([]byte, 0)
		for _, c := range line {
			row = append(row, byte(c))
		}

		out = append(out, row)
	}

	return out
}
