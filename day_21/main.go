package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Example")
	solvePart("./input/example.txt")

	fmt.Println("\nInput")
	solvePart("./input/input.txt")
}

func solvePart(location string) {
	fmt.Println("Part 1:", solve(location, 3))
	fmt.Println("Part 2:", solve(location, 26))
}

func solve(location string, depth int) int {
	file, _ := os.Open(location)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0

	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())
		if code == "" {
			continue
		}

		// Extract numeric value
		numStr := strings.TrimSuffix(code, "A")
		numStr = strings.TrimLeft(numStr, "0")
		numValue := 0
		if numStr != "" {
			numValue, _ = strconv.Atoi(numStr)
		}

		// Calculate sequence length
		seqLen := solveString(code, depth)
		total += seqLen * numValue
	}

	return total
}

type MovementPair struct {
	from, to string
}

var (
	locX = map[string]int{
		"7": 0, "8": 0, "9": 0,
		"4": 1, "5": 1, "6": 1,
		"1": 2, "2": 2, "3": 2,
		"#": 3, "^": 3, "A": 3,
		"<": 4, "v": 4, ">": 4,
	}
	locY = map[string]int{
		"7": 0, "8": 1, "9": 2,
		"4": 0, "5": 1, "6": 2,
		"1": 0, "2": 1, "3": 2,
		"#": 0, "^": 1, "A": 2,
		"<": 0, "v": 1, ">": 2,
	}
)

func generatePairs() []MovementPair {
	options := []string{"A", "^", "<", "v", ">"}
	var pairs []MovementPair
	for _, x := range options {
		for _, y := range options {
			pairs = append(pairs, MovementPair{x, y})
		}
	}
	return pairs
}

func findPairIndex(pairs []MovementPair, from, to string) int {
	for i, p := range pairs {
		if p.from == from && p.to == to {
			return i
		}
	}
	return -1
}

func bestPath(x1, y1, x2, y2 int) string {
	var result strings.Builder

	left := strings.Repeat("<", max(0, y1-y2))
	right := strings.Repeat(">", max(0, y2-y1))
	up := strings.Repeat("^", max(0, x1-x2))
	down := strings.Repeat("v", max(0, x2-x1))

	hashX := locX["#"]
	hashY := locY["#"]

	if hashX == min(x1, x2) && hashY == min(y1, y2) {
		result.WriteString(down)
		result.WriteString(right)
		result.WriteString(up)
		result.WriteString(left)
	} else if hashX == max(x1, x2) && hashY == min(y1, y2) {
		result.WriteString(up)
		result.WriteString(right)
		result.WriteString(down)
		result.WriteString(left)
	} else {
		result.WriteString(left)
		result.WriteString(down)
		result.WriteString(up)
		result.WriteString(right)
	}

	result.WriteString("A")
	return result.String()
}

func createMatrix(pairs []MovementPair) [][]int {
	n := len(pairs)
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	for i, src := range pairs {
		path := bestPath(locX[src.from], locY[src.from], locX[src.to], locY[src.to])
		prev := "A"
		for _, curr := range path {
			currStr := string(curr)
			idx := findPairIndex(pairs, prev, currStr)
			if idx >= 0 {
				matrix[i][idx]++
			}
			prev = currStr
		}
	}
	return matrix
}

func multiplyMatrix(a, b [][]int) [][]int {
	n := len(a)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0
			for k := 0; k < n; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	return result
}

func matrixPower(matrix [][]int, power int) [][]int {
	if power == 1 {
		return matrix
	}
	if power%2 == 0 {
		half := matrixPower(matrix, power/2)
		return multiplyMatrix(half, half)
	}
	return multiplyMatrix(matrix, matrixPower(matrix, power-1))
}

func solveString(code string, depth int) int {
	if depth == 0 {
		return len(code)
	}

	code = strings.ReplaceAll(code, "0", "^")
	pairs := generatePairs()
	matrix := createMatrix(pairs)

	powMatrix := matrixPower(matrix, depth-1)

	n := len(pairs)
	vector := make([]int, n)
	for i := range vector {
		vector[i] = 1
	}

	fastestPairs := make([]int, n)
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			sum += powMatrix[i][j] * vector[j]
		}
		fastestPairs[i] = sum
	}

	result := 0
	prev := "A"
	for _, curr := range code {
		currStr := string(curr)
		path := bestPath(locX[prev], locY[prev], locX[currStr], locY[currStr])

		pathPrev := "A"
		for _, move := range path {
			moveStr := string(move)
			idx := findPairIndex(pairs, pathPrev, moveStr)
			if idx >= 0 {
				result += fastestPairs[idx]
			}
			pathPrev = moveStr
		}
		prev = currStr
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
