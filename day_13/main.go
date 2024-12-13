package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Example")
	location := "./input/example.txt"
	machines := readInputFile(location)

	one, two := solve(machines, false), solve(machines, true)
	fmt.Println("Part One:", one)
	fmt.Println("Part Two:", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	machines = readInputFile(location)
	one, two = solve(machines, false), solve(machines, true)
	fmt.Println("Part One:", one)
	fmt.Println("Part Two:", two)
}

func solve(machines []Machine, partTwo bool) int {
	/* Equations
		* -------------------------------------------------
		* | x * buttonA.moveX + y * buttonB.moveX = clawX |
		* | x * buttonA.moveY + y * buttonB.moveY = clawY |
		* -------------------------------------------------
	* Linear Equations with two unknowns -> i can solve this still! :D
	*/
	var res int
	var coordinateMod int
	if partTwo {
		coordinateMod = 10000000000000
	}

	for _, machine := range machines {
		clawY := machine.clawY + coordinateMod
		clawX := machine.clawX + coordinateMod
		aMoveX := machine.buttonA.moveX
		aMoveY := machine.buttonA.moveY
		bMoveX := machine.buttonB.moveX
		bMoveY := machine.buttonB.moveY

		x, y := solveEquation(aMoveX, aMoveY, bMoveX, bMoveY, clawX, clawY)
		if x == -1 || y == -1 {
			continue
		}

		if !partTwo {
			if x > 100 || y > 100 {
				continue
			}
		}

		res += 3*x + y
	}

	return res
}

func solveEquation(aMoveX, aMoveY, bMoveX, bMoveY, clawX, clawY int) (int, int) {
	// Solving lineary equations is something i can still do it seems :)
	// It was a pain to get this working with the floating point precission
	p1 := clawY*bMoveX - bMoveY*clawX
	p2 := bMoveX*aMoveY - aMoveX*bMoveY

	// Gotta love floating points
	if !isInt(float64(p1) / float64(p2)) {
		return -1, -1
	}

	x := p1 / p2

	intermediate := clawX - x*aMoveX
	if !isInt(float64(intermediate) / float64(bMoveX)) {
		return -1, -1

	}

	y := intermediate / bMoveX

	if y < 0 || x < 0 {
		return -1, -1
	}

	return int(x), int(y)
}

func isInt(value float64) bool {
	return value == float64(int(value))
}

type Button struct {
	moveX int
	moveY int
}

type Machine struct {
	buttonA Button
	buttonB Button
	clawX   int
	clawY   int
}

func readInputFile(location string) []Machine {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	machines := []Machine{}
	var part int = 0
	var index int = 0
	var buttonA Button
	var buttonB Button
	var clawX int
	var clawY int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			part = 0
			index++
			continue
		}

		switch part {
		case 0:
			// buttonA
			buttonA = getButton(line)
		case 1:
			// buttonB
			buttonB = getButton(line)
		case 2:
			// claw
			clawX, clawY = getClaw(line)
		}

		part++

		if part == 3 {
			machine := Machine{
				buttonA: buttonA,
				buttonB: buttonB,
				clawX:   clawX,
				clawY:   clawY,
			}
			machines = append(machines, machine)
		}
	}

	return machines
}

func getButton(line string) Button {
	commaSplit := strings.Split(strings.Split(line, ":")[1], ",")
	numberOne := strings.Split(strings.Fields(commaSplit[0])[0], "+")[1]
	numberTwo := strings.Split(strings.Fields(commaSplit[1])[0], "+")[1]

	button := Button{}
	button.moveX, _ = strconv.Atoi(numberOne)
	button.moveY, _ = strconv.Atoi(numberTwo)
	return button
}

func getClaw(line string) (int, int) {
	commaSplit := strings.Split(strings.Split(line, ":")[1], ",")
	numberOne := strings.Split(strings.Fields(commaSplit[0])[0], "=")[1]
	numberTwo := strings.Split(strings.Fields(commaSplit[1])[0], "=")[1]

	clawX, _ := strconv.Atoi(numberOne)
	clawY, _ := strconv.Atoi(numberTwo)
	return clawX, clawY
}
