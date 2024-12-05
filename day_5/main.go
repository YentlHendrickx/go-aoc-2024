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

	one, two := solve(precedence, updates)
	fmt.Println("PartOne: ", one)
	fmt.Println("PartTwo: ", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"

	precedence, updates = readInputFile(location)

	one, two = solve(precedence, updates)
	fmt.Println("PartOne: ", one)
	fmt.Println("PartTwo: ", two)
}

func solve(prec Precedence, updates []Update) (int, int) {
	var one, two int = 0, 0

	correctUpdates, incorrectUpdates := getCorrectIncorrect(prec, updates)
	one = calculateMiddleSum(correctUpdates)

	two = calculateMiddleSum(correctUpdateSet(prec, incorrectUpdates))

	return one, two
}

func correctUpdateSet(prec Precedence, updates []Update) []Update {
	var output []Update
	for _, update := range updates {
		var currentSet Update = update
		for {
			correct, first, second := isCorrectUpdateSet(prec, currentSet)
			if correct {
				break
			}

			// Create new set which is the same as the current set but with the first and second swapped
			var newSet Update = make(Update, len(currentSet))
			for index, page := range currentSet {
				if page == first {
					newSet[index] = second
				} else if page == second {
					newSet[index] = first
				} else {
					newSet[index] = page
				}
			}

			currentSet = newSet
		}

		output = append(output, currentSet)
	}

	return output
}

func calculateMiddleSum(updates []Update) int {
	res := 0
	for _, update := range updates {
		middle := len(update) / 2
		res += update[middle]
	}

	return res
}

func getCorrectIncorrect(prec Precedence, updates []Update) ([]Update, []Update) {
	var correctUpdates []Update
	var incorrectUpdates []Update

	for _, update := range updates {
		correct, _, _ := isCorrectUpdateSet(prec, update)

		if correct {
			correctUpdates = append(correctUpdates, update)
		} else {
			incorrectUpdates = append(incorrectUpdates, update)
		}
	}

	return correctUpdates, incorrectUpdates
}

func isCorrectUpdateSet(prec Precedence, update Update) (bool, int, int) {
	correct := true
	var first, second int = -1, -1

	for index, page := range update {
		if _, ok := prec[page]; ok {
			curr := prec[page]
			for i := 0; i < index; i++ {
				if contains(curr, update[i]) {
					first = update[i]
					second = page
					correct = false
					break
				}
			}

		}
	}

	return correct, first, second
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
