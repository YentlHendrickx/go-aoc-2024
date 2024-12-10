package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	fmt.Println("\nExample")
	var location string = "./input/example.txt"

	memSections := readInputFile(location)

	fmt.Println("PartOne: ", solve(memSections, false))
	fmt.Println("PartTwo: ", solve(memSections, true))

	fmt.Println("\nInput")
	location = "./input/input.txt"

	memSections = readInputFile(location)

	fmt.Println("PartOne: ", solve(memSections, false))
	fmt.Println("PartTwo: ", solve(memSections, true))
}

const MUL_PAT = `mul\((\d+),(\d+)\)`
const DO_PAT = `do\(\)`
const DONT_PAT = `don't\(\)`

type Instruction struct {
	operation string
	one       int
	two       int
}

func solve(memSections []string, partTwo bool) int {
	var res int = 0

	instructions := make(map[int]Instruction)

	var indexOffset int = 0
	for _, section := range memSections {
		mulMatches := regexp.MustCompile(MUL_PAT).FindAllStringSubmatchIndex(section, -1)
		dos := regexp.MustCompile(DO_PAT).FindAllStringSubmatchIndex(section, -1)
		donts := regexp.MustCompile(DONT_PAT).FindAllStringSubmatchIndex(section, -1)

		if len(dos) > 0 {
			for _, do := range dos {
				instructions[indexOffset+do[0]] = Instruction{operation: "do", one: 0, two: 0}
			}
		}

		if len(donts) > 0 {
			for _, dont := range donts {
				instructions[indexOffset+dont[0]] = Instruction{operation: "dont", one: 0, two: 0}
			}
		}

		for _, mulIndex := range mulMatches {
			one := section[mulIndex[2]:mulIndex[3]]
			two := section[mulIndex[4]:mulIndex[5]]

			oneInt, _ := strconv.Atoi(one)
			twoInt, _ := strconv.Atoi(two)

			instructions[indexOffset+mulIndex[0]] = Instruction{operation: "mul", one: oneInt, two: twoInt}
		}

		indexOffset += len(section)
	}

	sortedSlice := slices.Sorted(maps.Keys(instructions))

	var do bool = true
	for _, index := range sortedSlice {
		switch instructions[index].operation {
		case "do":
			do = true
		case "dont":
			do = false
		case "mul":
			if do || !partTwo {
				res += instructions[index].one * instructions[index].two
			}
		}
	}

	return res
}

func readInputFile(location string) []string {
	file, err := os.Open(location)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result []string

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}
