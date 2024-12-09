package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
)

func main() {
	fmt.Println("Example")
	var location string = "./input/example.txt"
	diskMap := readInputFile(location)

	one := solve(diskMap)
	fmt.Println("PartOne: ", one)

	two := solve2(diskMap)
	fmt.Println("PartTwo: ", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	diskMap = readInputFile(location)

	one = solve(diskMap)
	fmt.Println("PartOne: ", one)

	two = solve2(diskMap)
	fmt.Println("PartTwo: ", two)
}

func solve(diskMap DiskMap) int {
	uncompressed := uncompressedMap(diskMap)
	squashed := squash(uncompressed)
	return calculateChecksun(squashed)
}

func solve2(diskMap DiskMap) int {
	uncompressed := uncompressedMap(diskMap)
	squashed := squash2(uncompressed)
	checksum := calculateChecksun(squashed)
	return checksum
}

func squash2(uncompressed []int) []int {
	uniqueSequences := getUniqueSequences(uncompressed)

	sortedSeq := slices.Sorted(maps.Keys(uniqueSequences))
	slices.Reverse(sortedSeq)

	for i := 0; i < len(sortedSeq); i++ {
		kS := sortedSeq[i]
		vS := uniqueSequences[kS]

		kIndex := -1
		for k := 0; k < len(uncompressed); k++ {
			if uncompressed[k] == kS {
				kIndex = k
				break
			}
		}

		emptySpace := getEmptySpace(uncompressed, kIndex)
		sortedEmpt := slices.Sorted(maps.Keys(emptySpace))

		for j := 0; j < len(sortedEmpt); j++ {
			eIndex := sortedEmpt[j]
			eLength := emptySpace[eIndex]
			if eLength < vS {
				continue
			}

			if eLength >= vS {
				// Use eIndex as the starting point, fill in uncompressed with kS
				for k := eIndex; k < eIndex+vS; k++ {
					uncompressed[k] = kS
				}

				// Remove the kS from the uncompressed -> end
				for k := eIndex + vS; k < len(uncompressed); k++ {
					if uncompressed[k] == kS {
						uncompressed[k] = -1
					}
				}
				break
			}
		}
	}

	return uncompressed
}

func getUniqueSequences(uncompressed []int) map[int]int {
	var uniqueSequences map[int]int = make(map[int]int)
	var currentSequence int = -1
	var length int = 0

	for i := 0; i < len(uncompressed)+1; i++ {
		if i == len(uncompressed) && currentSequence != -1 {
			uniqueSequences[currentSequence] = length
			break
		}

		if uncompressed[i] == -1 {
			continue
		}

		if currentSequence == -1 {
			currentSequence = uncompressed[i]
			length = 1
		} else if currentSequence != uncompressed[i] {
			uniqueSequences[currentSequence] = length
			currentSequence = uncompressed[i]
			length = 1
		} else {
			length++
		}
	}

	return uniqueSequences
}

func getEmptySpace(uncompressed []int, endIndex int) map[int]int {
	var emptySpace map[int]int = make(map[int]int)
	var space int = 0
	var currentStart int = -1
	for i := 0; i < endIndex; i++ {
		if uncompressed[i] == -1 {
			space++

			if currentStart == -1 {
				currentStart = i
			}

			emptySpace[currentStart] = space
		} else {
			space = 0
			currentStart = -1
		}
	}

	return emptySpace
}

func squash(uncompressed []int) []int {
	for i := 0; i < len(uncompressed); i++ {
		if uncompressed[i] != -1 {
			continue
		}

		for j := len(uncompressed) - 1; j > i; j-- {
			if uncompressed[j] != -1 {
				uncompressed[i] = uncompressed[j]
				uncompressed[j] = -1
				break
			}
		}
	}

	return uncompressed
}

func calculateChecksun(squashed []int) int {
	var res int = 0
	for i := 0; i < len(squashed); i++ {
		if squashed[i] == -1 {
			continue
		}

		res += squashed[i] * i
	}

	return res
}

func uncompressedMap(d DiskMap) []int {
	var id int = 0
	var result []int
	for i := 0; i < len(d); i += 2 {
		for j := 0; j < d[i]; j++ {
			result = append(result, id)
		}

		if i+1 >= len(d) {
			break
		}

		freeSpace := d[i+1]
		for j := 0; j < freeSpace; j++ {
			result = append(result, -1)
		}
		id++
	}

	return result
}

type DiskMap []int

func readInputFile(location string) DiskMap {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	out := make(DiskMap, 0)

	for scanner.Scan() {
		line := scanner.Text()
		for _, c := range line {
			out = append(out, int(c-'0'))
		}
	}

	return out
}
