package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("\nExample")
	var location string = "./input/example.txt"

	startMap, guard := readInputFile(location)
	one, path := solve(startMap, guard)
	fmt.Println("PartOne: ", one)

	two := solve2(startMap, guard, path)
	fmt.Println("PartTwo: ", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"

	startMap, guard = readInputFile(location)
	one, path = solve(startMap, guard)
	fmt.Println("PartOne: ", one)

	two = solve2(startMap, guard, path)
	fmt.Println("PartTwo: ", two)
}

func solve2(startMap Map, guard Guard, path []Guard) int {
	var result int = 0

	for _, guardPath := range path {
		var x = guardPath.x
		var y = guardPath.y

		if startMap[y][x] == 1 || guard.x == x && guard.y == y {
			continue
		}

		var newMap Map = make([][]Tile, len(startMap))
		for i := range startMap {
			newMap[i] = make([]Tile, len(startMap[i]))
			copy(newMap[i], startMap[i])
		}

		newMap[y][x] = 1

		locationHistory := make(map[int]map[int]map[int]bool)
		var newGuard = Guard{guard.x, guard.y, guard.direction}

		for {

			var newX, newY int = newGuard.x, newGuard.y

			switch newGuard.direction {
			case DIRECTION_UP:
				newY--
			case DIRECTION_RIGHT:
				newX++
			case DIRECTION_DOWN:
				newY++
			case DIRECTION_LEFT:
				newX--
			}

			if newX < 0 || newX >= len(startMap[0]) || newY < 0 || newY >= len(startMap) {
				break
			}

			if newMap[newY][newX] == 0 {
				newGuard.x = newX
				newGuard.y = newY
				continue
			}

			newGuard.direction = newGuard.direction + 1
			if newGuard.direction > 3 {
				newGuard.direction = 0
			}

			if locationHistory[newGuard.y][newGuard.x][newGuard.direction] == true {
				result++
				break
			}

			if locationHistory[newGuard.y] == nil {
				locationHistory[newGuard.y] = make(map[int]map[int]bool)
			}

			if locationHistory[newGuard.y][newGuard.x] == nil {
				locationHistory[newGuard.y][newGuard.x] = make(map[int]bool)
			}

			locationHistory[newGuard.y][newGuard.x][newGuard.direction] = true
		}
	}

	return result
}

func solve(startMap Map, guard Guard) (int, []Guard) {
	var result int = 0
	var visited [][]bool = make([][]bool, len(startMap))
	var path []Guard = make([]Guard, 0)

	for {
		if visited[guard.y] == nil {
			visited[guard.y] = make([]bool, len(startMap[0]))
		}

		if !visited[guard.y][guard.x] {
			path = append(path, guard)
			visited[guard.y][guard.x] = true
			result++
		}

		// Obstacle in front of the guard
		var newX, newY int = guard.x, guard.y

		switch guard.direction {
		case DIRECTION_UP:
			newY--
		case DIRECTION_RIGHT:
			newX++
		case DIRECTION_DOWN:
			newY++
		case DIRECTION_LEFT:
			newX--
		}

		if newX < 0 || newX >= len(startMap[0]) || newY < 0 || newY >= len(startMap) {
			break
		}

		// If obstacle at new position, turn right
		if startMap[newY][newX] == 1 {
			guard.direction = guard.direction + 1
			if guard.direction > 3 {
				guard.direction = 0
			}
		} else {
			guard.x = newX
			guard.y = newY
		}
	}

	return result, path
}

type Tile int
type Map [][]Tile

const DIRECTION_UP = 0
const DIRECTION_RIGHT = 1
const DIRECTION_DOWN = 2
const DIRECTION_LEFT = 3

type Guard struct {
	x, y      int
	direction int
}

func readInputFile(location string) (Map, Guard) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var outMap Map = make([][]Tile, 0)
	var guard Guard

	for scanner.Scan() {
		line := scanner.Text()
		var row []Tile = make([]Tile, 0)
		for _, char := range line {
			switch char {
			case '.':
				row = append(row, 0)
			case '#':
				row = append(row, 1)
			case '^':
				guard = Guard{len(row), len(outMap), DIRECTION_UP}
				row = append(row, 0)
			case '>':
				guard = Guard{len(row), len(outMap), DIRECTION_RIGHT}
				row = append(row, 0)
			case 'v':
				guard = Guard{len(row), len(outMap), DIRECTION_DOWN}
				row = append(row, 0)
			case '<':
				guard = Guard{len(row), len(outMap), DIRECTION_LEFT}
				row = append(row, 0)
			}
		}
		outMap = append(outMap, row)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return outMap, guard
}
