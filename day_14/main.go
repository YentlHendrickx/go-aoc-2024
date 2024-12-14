package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ROOM_X = 101
var ROOM_Y = 103

func main() {
	fmt.Println("\nInput")
	location := "./input/input.txt"
	robotMap := readInputFile(location)
	fmt.Println("Part One:", solve(robotMap))
	fmt.Println("Part Two:", solveTwo(robotMap))
}

func solve(robotMap RobotMap) int {
	for i := 0; i < 100; i++ {
		var roboCopy RobotMap = make(RobotMap)
		for location, robots := range robotMap {
			for _, robot := range robots {
				newX, newY := stepTime(location.x, location.y, robot.vx, robot.vy)
				newX, newY = limitInBounds(newX, newY)

				if _, ok := roboCopy[Location{newX, newY}]; ok {
					roboCopy[Location{newX, newY}] = append(roboCopy[Location{newX, newY}], robot)

				} else {
					roboCopy[Location{newX, newY}] = []Robot{robot}
				}
			}
		}

		robotMap = roboCopy
	}

	return countInQuadrant(robotMap)
}

func solveTwo(robotMap RobotMap) int {
	var timeStep int
	for {
		if checkRow(robotMap, timeStep) {
			break
		}
		timeStep++

		if timeStep > 100000 {
			break
		}

		var roboCopy RobotMap = make(RobotMap)
		for location, robots := range robotMap {
			for _, robot := range robots {
				newX, newY := stepTime(location.x, location.y, robot.vx, robot.vy)
				newX, newY = limitInBounds(newX, newY)

				if _, ok := roboCopy[Location{newX, newY}]; ok {
					roboCopy[Location{newX, newY}] = append(roboCopy[Location{newX, newY}], robot)

				} else {
					roboCopy[Location{newX, newY}] = []Robot{robot}
				}
			}
		}

		robotMap = roboCopy
	}

	return timeStep
}

func checkRow(robotMap RobotMap, steps int) bool {
	var uniqueCount int
	for location := range robotMap {
		if len(robotMap[location]) > 1 {
			break
		}

		uniqueCount++
	}

	// Won't lie found this by trial and error, since we don't know what the easter egg looks like this was the easiest way
	// I just increased uniqueCount > x until I found the right number -> turns out it's length of bots so each is at a unique location!
	if uniqueCount == len(robotMap) {
		printMap(robotMap, steps)
		return true
	}

	return false
}

func printMap(robotMap RobotMap, steps int) {
	fmt.Println("Steps:", steps)
	for y := 0; y < ROOM_Y; y++ {
		for x := 0; x < ROOM_X; x++ {
			if _, ok := robotMap[Location{x, y}]; ok {
				// fmt.Print(len(robotMap[Location{x, y}]))
				fmt.Print("█")
			} else {
				fmt.Print("░")
			}
		}
		fmt.Println()
	}
	fmt.Println()

	// Write this to a file if it doesn't exist create it
	file, err := os.OpenFile("./easter_egg.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Steps: %d\n", steps))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(2)
	}

	for y := 0; y < ROOM_Y; y++ {
		for x := 0; x < ROOM_X; x++ {
			if _, ok := robotMap[Location{x, y}]; ok {
				_, err = file.WriteString("█")
				if err != nil {
					fmt.Println("Error writing to file:", err)
					os.Exit(2)
				}
			} else {
				_, err = file.WriteString("░")
				if err != nil {
					fmt.Println("Error writing to file:", err)
					os.Exit(2)
				}
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(2)
		}
	}

	_, err = file.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(2)
	}
}

func countInQuadrant(robotMap RobotMap) int {
	var quadrantOne, quadrantTwo, quadrantThree, quadrantFour int

	var quadrantLimitX = ROOM_X / 2

	var quadrantLimitY = ROOM_Y / 2

	for location := range robotMap {
		if location.x < quadrantLimitX {
			if location.y < quadrantLimitY {
				quadrantOne += len(robotMap[location])
			} else if location.y > quadrantLimitY {
				quadrantThree += len(robotMap[location])
			}
		} else if location.x > quadrantLimitX {
			if location.y < quadrantLimitY {
				quadrantTwo += len(robotMap[location])
			} else if location.y > quadrantLimitY {
				quadrantFour += len(robotMap[location])
			}
		}
	}

	return quadrantOne * quadrantTwo * quadrantThree * quadrantFour
}

func limitInBounds(x, y int) (int, int) {
	if inBounds(x, y) {
		return x, y
	}

	for {
		if inBounds(x, y) {
			break
		}

		var xDiff, yDiff int
		if x < 0 {
			xDiff = -x
			x = ROOM_X - xDiff
		} else if x >= ROOM_X {
			xDiff = x - ROOM_X
			x = 0 + xDiff
		}

		if y < 0 {
			yDiff = -y
			y = ROOM_Y - yDiff
		} else if y >= ROOM_Y {
			yDiff = y - ROOM_Y
			y = 0 + yDiff
		}
	}

	return x, y
}

func inBounds(x, y int) bool {
	if x >= 0 && x < ROOM_X && y >= 0 && y < ROOM_Y {
		return true
	}

	return false
}

func stepTime(x, y int, vx, vy int) (int, int) {
	return x + vx, y + vy
}

type Robot struct {
	vx int
	vy int
}

type Location struct {
	x int
	y int
}

type RobotMap map[Location][]Robot

func readInputFile(location string) RobotMap {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	robotMap := make(RobotMap)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		position := strings.Split(parts[0], "=")[1]
		x, _ := strconv.Atoi(strings.Split(position, ",")[0])
		y, _ := strconv.Atoi(strings.Split(position, ",")[1])

		velocity := strings.Split(parts[1], "=")[1]
		vx, _ := strconv.Atoi(strings.Split(velocity, ",")[0])
		vy, _ := strconv.Atoi(strings.Split(velocity, ",")[1])

		if _, ok := robotMap[Location{x, y}]; ok {
			robotMap[Location{x, y}] = append(robotMap[Location{x, y}], Robot{vx, vy})
			continue
		}

		robotMap[Location{x, y}] = []Robot{{vx, vy}}
	}

	return robotMap
}
