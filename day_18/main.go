package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Example")
	exampleLocation := "./input/example.txt"
	coordinates := readInputFile(exampleLocation)
	one := len(solve(coordinates, 7, 7, 12))
	fmt.Println("Part One:", one)
	x, y := solve2(coordinates, 7, 7, 12)
	fmt.Print("Part Two: ")
	fmt.Print(x)
	fmt.Print(",")
	fmt.Println(y)

	fmt.Println("\nInput")
	inputLocation := "./input/input.txt"
	coordinates = readInputFile(inputLocation)
	one = len(solve(coordinates, 71, 71, 1024))
	fmt.Println("Part One:", one)
	x, y = solve2(coordinates, 71, 71, 1024)
	fmt.Print("Part Two: ")
	fmt.Print(x)
	fmt.Print(",")
	fmt.Println(y)
}

func solve2(coordinates []Coordinate, xSize, ySize, byteCount int) (int, int) {
	iterCount := 0

	previousPath := make(map[Coordinate]bool)
	cause := Coordinate{0, 0}
	for {
		cause = coordinates[int(iterCount+byteCount)-1]
		if _, found := previousPath[cause]; !found {
			if iterCount != 0 {
				iterCount++
				continue
			}
		}

		res := solve(coordinates, xSize, ySize, byteCount+iterCount)
		if res == nil {
			break
		}

		// Copy path to previousPath
		previousPath = make(map[Coordinate]bool)
		for k, v := range res {
			previousPath[k] = v
		}
		iterCount++
	}

	return cause.x, cause.y
}

func solve(coordinates []Coordinate, xSize, ySize, byteCount int) map[Coordinate]bool {
	grid := make([][]int, ySize)

	for i := 0; i < ySize; i++ {
		grid[i] = make([]int, xSize)
	}

	// Fill the grid with 0
	for i := 0; i < ySize; i++ {
		for j := 0; j < xSize; j++ {
			grid[i][j] = 0
		}
	}

	// Map coordinates to grid placing down the obstructions
	byteCounter := 0
	for _, c := range coordinates {
		if byteCounter == byteCount {
			break
		}

		grid[c.y][c.x] = 1
		byteCounter++
	}

	path := bfs(grid, xSize, ySize)
	return path
}

func printGrid(grid [][]int, xSize, ySize int, path map[Coordinate]bool) {
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			// Check if the current coordinate is in the path
			if _, found := path[Coordinate{x, y}]; found {
				fmt.Print("O")
			} else if grid[y][x] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}

		}
		fmt.Println()
	}
}

type Point struct {
	current Coordinate
	path    map[Coordinate]bool
}

func bfs(grid [][]int, xSize, ySize int) map[Coordinate]bool {
	// Find the shortest path, starting from top left 0,0 to bottom right xSize, ySize. BFS
	var stack []Point = []Point{{Coordinate{0, 0}, make(map[Coordinate]bool)}}
	visited := make(map[Coordinate]bool)

	for len(stack) > 0 {
		point := stack[0]
		stack = stack[1:]

		current := point.current

		// Check if we have already visited this coordinates
		if _, ok := visited[current]; ok {
			continue
		}

		visited[current] = true

		if current.x == xSize-1 && current.y == ySize-1 {
			return point.path
		}

		directions := []Coordinate{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
		for _, d := range directions {
			if !(current.x+d.x >= 0 && current.x+d.x < xSize && current.y+d.y >= 0 && current.y+d.y < ySize) {
				continue
			}

			if grid[current.y+d.y][current.x+d.x] == 1 {
				continue
			}

			// Copy path
			newPath := make(map[Coordinate]bool)
			for k, v := range point.path {
				newPath[k] = v
			}
			newPath[current] = true

			// Add the current point to the path
			stack = append(stack, Point{Coordinate{current.x + d.x, current.y + d.y}, newPath})
		}
	}

	return nil
}

type Coordinate struct {
	x int
	y int
}

func readInputFile(location string) []Coordinate {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	coordinates := []Coordinate{}

	for scanner.Scan() {
		line := scanner.Text()
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)

		coordinates = append(coordinates, Coordinate{x, y})
	}

	return coordinates
}
