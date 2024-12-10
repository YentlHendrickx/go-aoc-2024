package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Example")
	var location string = "./input/example.txt"
	grid := readInputFile(location)

	one, two := solve(grid)
	fmt.Println("PartOne: ", one)
	fmt.Println("PartTwo: ", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	grid = readInputFile(location)

	one, two = solve(grid)
	fmt.Println("PartOne: ", one)
	fmt.Println("PartTwo: ", two)
}

type Stack struct {
	x          int
	y          int
	trailheadX int
	trailheadY int
}

func solve(grid Grid) (int, int) {
	var stack []Stack = make([]Stack, 0)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 0 {
				stack = append(stack, Stack{i, j, i, j})
			}
		}
	}

	var visitedEnd map[int]map[int][][]int = make(map[int]map[int][][]int)

	var resTwo int = 0
	for len(stack) > 0 {
		var currentEntry Stack = stack[0]
		currentX, currentY, trailHeadX, trailHeadY := currentEntry.x, currentEntry.y, currentEntry.trailheadX, currentEntry.trailheadY
		stack = stack[1:]

		currentValue := grid[currentX][currentY]

		xRange := []int{-1, 0, 0, 1}
		yRange := []int{0, 1, -1, 0}
		for i := 0; i < 4; i++ {
			newX := currentX + xRange[i]
			newY := currentY + yRange[i]

			if newX < 0 || newX >= len(grid) || newY < 0 || newY >= len(grid[0]) {
				continue
			}

			if grid[newX][newY] != currentValue+1 {
				continue
			}

			if currentValue+1 != 9 {
				stack = append(stack, Stack{newX, newY, currentEntry.trailheadX, currentEntry.trailheadY})
				continue
			}

			if _, ok := visitedEnd[trailHeadX]; !ok {
				visitedEnd[trailHeadX] = make(map[int][][]int)
			}

			if _, ok := visitedEnd[trailHeadX][trailHeadY]; !ok {
				visitedEnd[trailHeadX][trailHeadY] = make([][]int, 0)
			}

			// Check if x,y is already in the visitedEnd
			var found bool = false
			for _, entry := range visitedEnd[trailHeadX][trailHeadY] {
				if entry[0] == newX && entry[1] == newY {
					found = true
				}
			}

			if !found {
				visitedEnd[trailHeadX][trailHeadY] = append(visitedEnd[trailHeadX][trailHeadY], []int{newX, newY})
			}

			resTwo++
		}
	}

	var resOne int = 0

	for trailX, entry := range visitedEnd {
		for trailY := range entry {
			for i := 0; i < len(visitedEnd[trailX][trailY]); i++ {
				resOne++
			}
		}
	}

	return resOne, resTwo
}

type Grid [][]int

func readInputFile(location string) Grid {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	out := make(Grid, 0)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, 0)
		for _, r := range line {
			int, _ := strconv.Atoi(string(r))
			row = append(row, int)
		}

		out = append(out, row)
	}

	return out
}
