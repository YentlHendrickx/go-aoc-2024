package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Example")
	solvePart("./input/example.txt")

	fmt.Println("\nInput")
	solvePart("./input/input.txt")
}

func solvePart(location string) {
	one, two := solve(location)
	fmt.Println("Part 1:", one)
	fmt.Println("Part 2:", two)
}

type Consecutive struct {
	n1 int
	n2 int
	n3 int
	n4 int
}

func solve(location string) (int, int) {
	nums := readInputFile(location)
	resOne := 0

	var sequences [][]int = make([][]int, len(nums))

	for index, num := range nums {
		inter := num
		for i := 0; i < 2000; i++ {
			inter = evolve(inter)
			sequences[index] = append(sequences[index], inter%10)
		}

		resOne += inter
	}

	var changes []map[Consecutive]int = make([]map[Consecutive]int, len(nums))
	for index, sequence := range sequences {
		changes[index] = make(map[Consecutive]int)
		for i := 1; i < len(sequence)-3; i++ {
			consecutive := Consecutive{sequence[i] - sequence[i-1], sequence[i+1] - sequence[i], sequence[i+2] - sequence[i+1], sequence[i+3] - sequence[i+2]}

			if _, found := changes[index][consecutive]; !found {
				changes[index][consecutive] = sequence[i+3]
			}
		}
	}

	// Build up each unique 4 number sequence that is possible, a sequence can be from -9 to 9 and is 4 numbers
	// This is a brute force solution, but it works for the given input
	var possibleSequences []Consecutive
	for i := -9; i <= 9; i++ {
		for j := -9; j <= 9; j++ {
			for k := -9; k <= 9; k++ {
				for l := -9; l <= 9; l++ {
					// If all 4 are negative, skip
					if i < 0 && j < 0 && k < 0 && l < 0 {
						continue
					}

					possibleSequences = append(possibleSequences, Consecutive{i, j, k, l})
				}
			}
		}
	}

	var maxRes int
	for _, sequence := range possibleSequences {
		var res int
		for index, change := range changes {
			// Impossible solution
			if res+(len(changes)-index)*9 < maxRes {
				break
			}

			if val, found := change[sequence]; found {
				res += val
			}
		}

		if res > maxRes {
			maxRes = res
		}
	}
	resTwo := maxRes

	return resOne, resTwo
}

func evolve(number int) int {
	// Step one
	stepOne := prune(mix(number*64, number))

	stepOneFloat := float64(stepOne)
	stepTwo := int(math.Floor(stepOneFloat / 32))
	stepTwo = prune(mix(stepTwo, stepOne))

	stepThree := prune(mix(stepTwo*2048, stepTwo))

	return stepThree
}

func mix(n1, n2 int) int {
	return n1 ^ n2
}

func prune(number int) int {
	return number % 16777216
}

func readInputFile(location string) []int {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res []int
	for scanner.Scan() {
		line := scanner.Text()
		intVal, _ := strconv.Atoi(line)
		res = append(res, intVal)
	}

	return res
}
