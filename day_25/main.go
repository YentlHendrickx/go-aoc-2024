package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Example")
	fmt.Println("Part One:", solve("./input/example.txt"))

	fmt.Println("\nInput")
	fmt.Println("Part One:", solve("./input/input.txt"))

	// And with that final print statement, we are don!
	// I had a lot of fun doing this year's AOC, some challenges were really hard, but I learned a lot.
	// Eric Wastl is a genius. <3
}

func solve(location string) int {
	blocks, _ := readInputFile(location)

	locks := make([]Block, 0)
	keys := make([]Block, 0)
	for _, block := range blocks {
		if block.IsLock {
			locks = append(locks, block)
			continue
		}

		keys = append(keys, block)
	}

	var result int
	for _, lock := range locks {
		for _, key := range keys {
			if isMatch(lock.Heights, key.Heights) {
				result += 1
			}
		}
	}

	return result
}

func isMatch(lock Sequence, key Sequence) bool {
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] >= 6 {
			return false
		}
	}

	return true
}

type Sequence []int

type Block struct {
	IsLock  bool
	Grid    []string
	Heights Sequence
}

func readInputFile(location string) ([]Block, error) {
	file, err := os.Open(location)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var blocks []Block
	var currentBlock Block

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			if len(currentBlock.Grid) > 0 {
				currentBlock.Heights = computeHeights(currentBlock.Grid, currentBlock.IsLock)
				blocks = append(blocks, currentBlock)
				currentBlock = Block{}
			}
			continue
		}

		if len(currentBlock.Grid) == 0 {
			currentBlock.IsLock = strings.HasPrefix(line, "#")
		}

		currentBlock.Grid = append(currentBlock.Grid, line)
	}

	if len(currentBlock.Grid) > 0 {
		currentBlock.Heights = computeHeights(currentBlock.Grid, currentBlock.IsLock)
		blocks = append(blocks, currentBlock)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return blocks, nil
}

func computeHeights(grid []string, isLock bool) Sequence {
	if len(grid) == 0 {
		return nil
	}

	numRows := len(grid)
	numCols := len(grid[0])
	heights := make(Sequence, numCols)

	for col := 0; col < numCols; col++ {
		for row := 0; row < numRows; row++ {
			if grid[row][col] == '#' && !isLock {
				heights[col] = numRows - row - 1
				break
			} else if grid[row][col] == '.' && isLock {
				heights[col] = row - 1
				break
			}
		}
	}

	return heights
}

func printBlocks(blocks []Block) {
	for i, block := range blocks {
		if block.IsLock {
			fmt.Printf("Lock #%d has pin heights %v:\n\n", i+1, block.Heights)
		} else {
			fmt.Printf("Key #%d heights %v:\n\n", i+1, block.Heights)
		}

		// Print the grid
		for _, line := range block.Grid {
			fmt.Println(line)
		}
		fmt.Println()
	}
}
