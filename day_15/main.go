package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {
	fmt.Println("Example")
	location := "./input/example_2.txt"
	wh, robot := readInputFile(location)
	whP2, robotP2 := createPartTwoMap(wh, robot)
	fmt.Println("Part One:", solve(wh, robot, false))
	fmt.Println("Part Two:", solve(whP2, robotP2, true))

	fmt.Println("\nInput")
	location = "./input/input.txt"
	wh, robot = readInputFile(location)
	whP2, robotP2 = createPartTwoMap(wh, robot)
	fmt.Println("Part One:", solve(wh, robot, false))
	fmt.Println("Part Two:", solve(whP2, robotP2, true))
}

func createPartTwoMap(wh Warehouse, robot Robot) (Warehouse, Robot) {
	var xValues []int
	var yValues []int
	for k := range wh {
		if !slices.Contains(xValues, k.x) {

			xValues = append(xValues, k.x)
		}

		if !slices.Contains(yValues, k.y) {
			yValues = append(yValues, k.y)
		}
	}
	// Since everything is doubled on x, robot position is also doubled
	robot.position.x *= 2

	slices.Sort(xValues)
	slices.Sort(yValues)

	var newWh Warehouse = make(Warehouse)
	for _, y := range yValues {
		newX := 1
		for _, x := range xValues {
			if wh[Coordinate{x, y}] == 'O' {
				newWh[Coordinate{newX - 1, y}] = '['
				newWh[Coordinate{newX, y}] = ']'
			} else {
				newWh[Coordinate{newX - 1, y}] = wh[Coordinate{x, y}]
				newWh[Coordinate{newX, y}] = wh[Coordinate{x, y}]
			}
			newX += 2
		}
	}

	return newWh, robot
}

func getNewPosition(currentPosition Coordinate, direction rune) Coordinate {
	newPosition := currentPosition
	switch direction {
	case '^':
		newPosition.y--
	case 'v':
		newPosition.y++
	case '>':
		newPosition.x++
	case '<':
		newPosition.x--
	}

	return newPosition
}

func getBoxOneBoxTwo(wh Warehouse, newPosition Coordinate, instruction rune) (Coordinate, Coordinate) {
	var boxPartOne Coordinate
	var boxPartTwo Coordinate
	if instruction == '<' {
		boxPartOne = getNewPosition(newPosition, instruction)
		boxPartTwo = newPosition
		return boxPartOne, boxPartTwo
	}

	if instruction == '>' {
		boxPartOne = newPosition
		boxPartTwo = getNewPosition(newPosition, instruction)
		return boxPartOne, boxPartTwo
	}

	if wh[newPosition] == '[' {
		boxPartOne = newPosition
		boxPartTwo = getNewPosition(newPosition, '>')
		return boxPartOne, boxPartTwo
	}

	boxPartTwo = newPosition
	boxPartOne = getNewPosition(newPosition, '<')
	return boxPartOne, boxPartTwo
}

func pushDoubleBox(wh Warehouse, boxPartOne Coordinate, boxPartTwo Coordinate, direction rune) (Warehouse, bool) {
	newOne := getNewPosition(boxPartOne, direction)
	newTwo := getNewPosition(boxPartTwo, direction)

	if checkValidPair(wh, newOne, newTwo) {
		wh[boxPartOne] = '.'
		wh[boxPartTwo] = '.'
		wh[newOne] = '['
		wh[newTwo] = ']'
		return wh, true
	} else if wh[newOne] == '#' || wh[newTwo] == '#' {
		return wh, false
	}

	// Deep copy the warehouse
	ogWarehouse := make(Warehouse)
	for k, v := range wh {
		ogWarehouse[k] = v
	}

	valid := true
	if (direction == '^' || direction == 'v') && (wh[newOne] == ']' || wh[newTwo] == '[') {
		nextState := wh
		if wh[newOne] == ']' {
			oneTopLeft := getNewPosition(newOne, '<')
			oneTopRight := getNewPosition(newTwo, '<')
			whStateTopLeft, topLeftResult := pushDoubleBox(wh, oneTopLeft, oneTopRight, direction)

			if !topLeftResult {
				return ogWarehouse, false
			}

			nextState = whStateTopLeft
		}

		if wh[newTwo] == '[' {
			twoTopLeft := getNewPosition(newOne, '>')
			twoTopRight := getNewPosition(newTwo, '>')

			whStateTopRight, topRightResult := pushDoubleBox(nextState, twoTopLeft, twoTopRight, direction) // use state top left, we revert if either result is false

			if !topRightResult {
				return ogWarehouse, false
			}

			nextState = whStateTopRight
		}

		wh = nextState
		valid = true
	}

	// if wh[newOne] == wh[boxPartOne] && wh[newTwo] == wh[boxPartTwo] {
	if ((wh[newOne] == '[' && wh[newTwo] == ']') && (direction == '^' || direction == 'v')) || (wh[newOne] == ']' && direction == '<') || (wh[newTwo] == '[' && direction == '>') {
		var nextOne, nextTwo Coordinate
		if direction == '^' || direction == 'v' {
			nextOne = newOne
			nextTwo = newTwo
		} else {
			nextOne = getNewPosition(newOne, direction)
			nextTwo = getNewPosition(newTwo, direction)
		}

		var nextState Warehouse
		nextState, valid = pushDoubleBox(wh, nextOne, nextTwo, direction)

		if valid {
			wh = nextState
		}
	}

	if valid {
		wh[boxPartOne] = '.'
		wh[boxPartTwo] = '.'
		wh[newOne] = '['
		wh[newTwo] = ']'
		return wh, true
	}

	return ogWarehouse, false
}

func checkValidPair(wh Warehouse, one Coordinate, two Coordinate) bool {
	if wh[one] == '#' || wh[two] == '#' {
		return false
	}

	if wh[one] == '.' && wh[two] == '.' {
		return true
	}

	return false
}

func solve(wh Warehouse, robot Robot, partTwo bool) int {
	for _, direction := range robot.instructions {
		newPosition := getNewPosition(robot.position, direction)

		if wh[newPosition] == '.' {
			robot.position = newPosition
			continue
		}

		if wh[newPosition] == '#' {
			continue
		}

		if partTwo {
			boxPartOne, boxPartTwo := getBoxOneBoxTwo(wh, newPosition, direction)
			var res bool
			wh, res = pushDoubleBox(wh, boxPartOne, boxPartTwo, direction)
			if res {
				robot.position = newPosition
			}

		} else {
			var res bool
			wh, res = pushBox(wh, newPosition, direction)
			if res {
				robot.position = newPosition
			}
		}

	}

	printWarehouse(wh, robot, ' ')
	return calculateCoordinates(wh, partTwo)
}

func calculateCoordinates(wh Warehouse, partTwo bool) int {
	var res int
	for k, v := range wh {
		if partTwo && v != '[' {
			continue
		} else if !partTwo && v != 'O' {
			continue
		}

		res += k.x + 100*k.y
	}

	return res
}

func pushBox(wh Warehouse, currentBox Coordinate, direction rune) (Warehouse, bool) {
	newPosition := getNewPosition(currentBox, direction)

	if wh[newPosition] == '.' {
		wh[newPosition] = 'O'
		wh[currentBox] = '.'
		return wh, true
	} else if wh[newPosition] == 'O' {
		newWh, res := pushBox(wh, newPosition, direction)
		if res {
			newWh[newPosition] = 'O'
			newWh[currentBox] = '.'
			return newWh, true
		}
	}

	return wh, false
}

func printWarehouse(wh Warehouse, robot Robot, direction rune) {
	var xValues []int
	var yValues []int
	for k := range wh {
		if !slices.Contains(xValues, k.x) {

			xValues = append(xValues, k.x)
		}

		if !slices.Contains(yValues, k.y) {
			yValues = append(yValues, k.y)
		}
	}

	slices.Sort(xValues)
	slices.Sort(yValues)

	for _, y := range yValues {
		for _, x := range xValues {
			if robot.position.x == x && robot.position.y == y {
				fmt.Print("@")
			} else {
				fmt.Print(string(wh[Coordinate{x, y}]))
			}
		}
		fmt.Println()
	}
	if direction != ' ' {
		fmt.Println("Direction:", string(direction))
	}

	fmt.Println()
}

type Coordinate struct {
	x int
	y int
}

type Warehouse map[Coordinate]rune
type Robot struct {
	position     Coordinate
	instructions []rune
}

func readInputFile(location string) (Warehouse, Robot) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	warehouse := make(Warehouse)
	robot := Robot{}

	var hasWarehouse bool
	var y int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			hasWarehouse = true
		}

		// Loop over chars in line, add to warehouse
		for i, char := range line {
			if !hasWarehouse {
				if char == '@' {
					robot.position = Coordinate{i, y}
					warehouse[Coordinate{i, y}] = '.'
				} else {
					warehouse[Coordinate{i, y}] = char
				}
			} else {
				robot.instructions = append(robot.instructions, char)
			}
		}

		y++
	}

	return warehouse, robot
}
