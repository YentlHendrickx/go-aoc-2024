// After benchmarking 1000 runs
// PartOne: ~460us
// ParTwo: ~15ms
// I can go up to blink 1000 in <500ms... SPEED
// I'm happy with the result, I'm sure there are more optimizations to be done, but I'm happy with the result.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Example")
	var location string = "./input/example.txt"
	nums := readInputFile(location)
	solve(nums)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	nums = readInputFile(location)
	solve(nums)
}

const PART_ONE_BLINKS = 25
const PART_TWO_BLINKS = 75

func solve(in []int) {
	var one int = 0
	var two int = 0

	startTime := time.Now()

	// Double caching approach, result cache -> num = 25 -> cache knows that result will be [2, 5].
	var resultCache map[int][]int = make(map[int][]int)

	// Duplicate cache, processing all numbers takes way too long
	// 3, 3, 4 -> 3: 2, 4: 1
	var duplicateCache map[int]int = make(map[int]int)

	for _, num := range in {
		if _, ok := duplicateCache[num]; ok {
			duplicateCache[num]++
			continue
		}

		duplicateCache[num] = 1
	}

	for i := 1; i <= PART_TWO_BLINKS; i++ {
		var newDuplicateCache map[int]int = make(map[int]int)
		for num, occ := range duplicateCache {
			blink := blink(num, &resultCache)

			for _, b := range blink {
				if _, ok := newDuplicateCache[b]; ok {
					newDuplicateCache[b] += occ
				} else {
					newDuplicateCache[b] = occ
				}
			}
		}

		duplicateCache = newDuplicateCache

		if i == PART_ONE_BLINKS {
			for _, occ := range duplicateCache {
				one += occ
			}

			elapsed := time.Since(startTime)
			fmt.Println("PartOne: ", one, "Time: ", elapsed)
		}
	}

	for _, occ := range duplicateCache {
		two += occ
	}

	elapsed := time.Since(startTime)
	fmt.Println("PartTwo: ", two, "Time: ", elapsed/1000)
}

func blink(num int, cache *map[int][]int) []int {
	if _, ok := (*cache)[num]; ok {
		return (*cache)[num]
	}

	if num == 0 {
		return []int{1}
	}

	asStr := strconv.Itoa(num)
	if len(asStr)%2 == 0 {
		numOne, _ := strconv.Atoi(asStr[:len(asStr)/2])
		numTwo, _ := strconv.Atoi(asStr[len(asStr)/2:])

		(*cache)[num] = []int{numOne, numTwo}
		return []int{numOne, numTwo}
	}

	newNum := num * 2024
	(*cache)[num] = []int{newNum}
	return []int{newNum}
}

func readInputFile(location string) []int {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	out := make([]int, 0)

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.Split(line, " ")
		for _, num := range split {
			i, _ := strconv.Atoi(num)
			out = append(out, i)
		}
	}

	return out
}
