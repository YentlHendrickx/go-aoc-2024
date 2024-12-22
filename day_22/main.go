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
	var memo map[int]int = make(map[int]int)
	one, two := solve(location, memo)
	fmt.Println("Part 1:", one)
	fmt.Println("Part 2:", two)
}

type Consecutive struct {
	n1 int
	n2 int
	n3 int
	n4 int
}

func solve(location string, memo map[int]int) (int, int) {
	nums := readInputFile(location)
	resOne := 0

	var sequences [][]int = make([][]int, len(nums))

	for index, num := range nums {
		inter := num
		for i := 0; i < 2000; i++ {
			if _, found := memo[inter]; found {
				inter = memo[inter]
			} else {
				sol := evolve(inter)
				memo[inter] = sol
				inter = sol
			}

			sequences[index] = append(sequences[index], inter%10)
		}

		resOne += inter
	}

	var changes []map[Consecutive]int = make([]map[Consecutive]int, len(nums))
	var possibleSequences map[Consecutive]int = make(map[Consecutive]int)
	var maxRes int
	for index, sequence := range sequences {
		changes[index] = make(map[Consecutive]int)
		for i := 1; i < len(sequence)-3; i++ {
			consecutive := Consecutive{sequence[i] - sequence[i-1], sequence[i+1] - sequence[i], sequence[i+2] - sequence[i+1], sequence[i+3] - sequence[i+2]}

			if _, found := changes[index][consecutive]; !found {
				changes[index][consecutive] = sequence[i+3]

				if _, found := possibleSequences[consecutive]; !found {
					possibleSequences[consecutive] = sequence[i+3]
				} else {
					newSeq := possibleSequences[consecutive] + sequence[i+3]
					if newSeq > maxRes {
						maxRes = newSeq
					}
					possibleSequences[consecutive] = newSeq
				}
			}
		}
	}

	return resOne, maxRes
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
