package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("\nExample")
	var location string = "./input/example.txt"

	reports := readInputFile(location)

	fmt.Println("Part one: ", solve(reports, false))
	fmt.Println("Part two: ", solve(reports, true))

	fmt.Println("\nInput")
	location = "./input/input.txt"

	reports = readInputFile(location)

	fmt.Println("Part one: ", solve(reports, false))
	fmt.Println("Part two: ", solve(reports, true))
}

func solve(reports [][]int, partTwo bool) int {
	count := 0
	for _, report := range reports {
		if partTwo {
			if isSafePartTwo(report) {
				count++
			}
		} else {
			if isSafe(report) {
				count++
			}
		}
	}

	return count
}

func isSafe(report []int) bool {
	dir := getDirection(report[0] - report[1])

	for i := 0; i < len(report)-1; i++ {
		diff := report[i] - report[i+1]
		newDir := getDirection(diff)

		if newDir != dir {
			return false
		}

		absDiff := int(math.Abs(float64(diff)))

		if absDiff > 3 || absDiff < 1 {
			return false
		}
	}

	return true
}

func isSafePartTwo(report []int) bool {
	// Basically we reslice the array, and check if it's safe
	for i := 0; i < len(report); i++ {
		var sliced []int
		for j := 0; j < len(report); j++ {
			if j == i {
				continue
			}

			sliced = append(sliced, report[j])
		}

		if isSafe(sliced) {
			return true
		}
	}

	return false
}

func getDirection(diff int) int {
	switch {
	case diff < 0:
		return 1
	case diff > 0:
		return 2
	default:
		return 0
	}
}

func readInputFile(location string) [][]int {
	file, err := os.Open(location)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}

	// Close it!
	defer file.Close()

	// Loop over all lines
	scanner := bufio.NewScanner(file)
	var result [][]int

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")

		var temp []int
		for _, i := range line {
			j, _ := strconv.Atoi(i)
			temp = append(temp, j)
		}

		result = append(result, temp)
	}

	return result
}
