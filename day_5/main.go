package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("\nExample")
	var location string = "./input/example.txt"

	precedence, updates := readInputFile(location)
	// fmt.Println("Precedence: ", precedence)
	// fmt.Println("\nUpdates: ", updates)

	fmt.Println("PartOne: ", partOne(precedence, updates))
	// // fmt.Println("PartTwo: ", partTwo(grid))
	//
	fmt.Println("\nInput")
	location = "./input/input.txt"

	precedence, updates = readInputFile(location)
	// fmt.Println("Precedence: ", precedence)
	// fmt.Println("\nUpdates: ", updates)

	fmt.Println("PartOne: ", partOne(precedence, updates))
}

func partOne(prec Precedence, updates []Update) int {
	var correctUpdates []Update
	for _, update := range updates {
		correct := true
		for index, page := range update {
			if _, ok := prec[page]; ok {
				curr := prec[page]
				for i := 0; i < index; i++ {
					if contains(curr, update[i]) {
						correct = false
						break
					}
				}

			}
		}

		if correct {
			correctUpdates = append(correctUpdates, update)
		}
	}

	res := 0

	for _, update := range correctUpdates {
		middle := len(update) / 2
		res += update[middle]
	}

	return res
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type Precedence map[int][]int

type Update []int

func readInputFile(location string) (Precedence, []Update) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var precedence Precedence = make(map[int][]int)
	var updates []Update

	var processedPairs bool = false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Check if we reached the end of pairs
		if line == "" {
			processedPairs = true
			continue
		}

		if !processedPairs {
			// Read pairs formatted as `int|int`
			parts := strings.Split(line, "|")

			// We have two parts, store the first part as the key and the second part as the value, if the key already exists append the value to the key
			firstPart, _ := strconv.Atoi(parts[0])
			secondPart, _ := strconv.Atoi(parts[1])

			if mapValue, ok := precedence[firstPart]; ok {
				precedence[firstPart] = append(mapValue, secondPart)
			} else {
				precedence[firstPart] = []int{secondPart}
			}

		} else {
			// Read updates formatted as `int,int,int,...`
			parts := strings.Split(line, ",")
			var update Update
			for _, part := range parts {
				page, err := strconv.Atoi(strings.TrimSpace(part))
				if err != nil {
					fmt.Println("Error parsing update value:", err)
					continue
				}
				update = append(update, page)
			}
			updates = append(updates, update)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return precedence, updates
}
