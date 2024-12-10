package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var location string = "./input/input.txt"
	var a, b = readInputFile(location)

	if a == nil || b == nil {
		fmt.Println("Error reading file")
		return
	}

	var res int = partOne(a, b)
	fmt.Println("Part One: ", res)

	res = partTwo(a, b)
	fmt.Println("Part Two: ", res)
}

func partOne(a []int, b []int) int {
	a = bubbleSort(a)
	b = bubbleSort(b)

	var res int = 0

	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		if diff < 0 {
			diff = diff * -1
		}

		res += diff
	}

	return res
}

func partTwo(a []int, b []int) int {
	var res int = 0

	var occMap map[int]int = make(map[int]int)

	for i := 0; i < len(a); i++ {
		curr := a[i]

		if occMap[curr] == 0 {
			occurences := countOccurences(b, a[i])
			occMap[curr] = occurences
			curr = occurences
		}

		res += curr * a[i]
	}

	return res
}

func countOccurences(haystack []int, needle int) int {
	var count int = 0

	for i := 0; i < len(haystack); i++ {
		if haystack[i] == needle {
			count++
		}
	}

	return count
}

func bubbleSort(a []int) []int {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a)-1; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
			}
		}
	}

	return a
}

func readInputFile(location string) ([]int, []int) {
	file, err := os.Open(location)

	if err != nil {
		fmt.Println("Error: ", err)
		return nil, nil
	}

	// Close it!
	defer file.Close()

	// Loop over all lines
	scanner := bufio.NewScanner(file)
	var a, b []int

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "   ")
		aStr, _ := strconv.Atoi(line[0])
		bStr, _ := strconv.Atoi(line[1])

		a = append(a, aStr)
		b = append(b, bStr)
	}

	return a, b
}
