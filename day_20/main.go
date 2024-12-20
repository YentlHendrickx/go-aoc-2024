package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type point struct {
	x, y int
}

type offset struct {
	point    point
	distance int
}

func main() {
	fmt.Println("Example")
	solveBoth("./input/example.txt", 50)

	fmt.Println("\nInput")
	solveBoth("./input/input.txt", 100)
}

func solveBoth(location string, savingThreshold int) {
	walls, sizeX, sizeY, start, end := readInputFile(location)
	one := solve(walls, sizeX, sizeY, start, end, 2, savingThreshold)
	fmt.Println("Part One: ", one)

	two := solve(walls, sizeX, sizeY, start, end, 20, savingThreshold)
	fmt.Println("Part Two: ", two)
}

func solve(walls map[point]bool, sizeX, sizeY int, start, end point, cheatSteps int, savingThreshold int) int {
	route := bfs(start, end, walls, sizeX, sizeY)
	if _, found := route[end]; !found {
		fmt.Println("No path found from Start to End.")
		return -1
	}

	savings := findShortcuts(route, cheatSteps)

	count := 0
	for saving := range savings {
		if saving >= savingThreshold {
			count += savings[saving]
		}
	}

	return count
}

func bfs(start, end point, walls map[point]bool, width, height int) map[point]int {
	queue := []point{start}
	visited := make(map[point]int)
	visited[start] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			return visited
		}

		for _, offset := range getOffsetsFromPoint(current, 1, width, height) {
			if _, found := visited[offset.point]; found {
				continue
			}
			if walls[offset.point] {
				continue
			}
			visited[offset.point] = visited[current] + 1
			queue = append(queue, offset.point)
		}
	}

	return visited
}

func findShortcuts(route map[point]int, radius int) map[int]int {
	shortcuts := make(map[int]int)
	for current, step := range route {
		offsets := getOffsetsFromPoint(current, radius, math.MaxInt32, math.MaxInt32)
		for _, off := range offsets {
			routeStep, inRoute := route[off.point]
			if inRoute {
				saving := routeStep - step - off.distance
				if saving > 0 {
					shortcuts[saving]++
				}
			}
		}
	}

	return shortcuts
}

func getOffsetsFromPoint(from point, radius, width, height int) []offset {
	var result []offset
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			candidate := point{from.x + dx, from.y + dy}

			if candidate.x < 0 || candidate.x >= width || candidate.y < 0 || candidate.y >= height {
				continue
			}

			dist := getDistance(from, candidate)
			if dist > 0 && dist <= radius {
				result = append(result, offset{candidate, dist})
			}
		}
	}
	return result
}

// Manhattan distance
func getDistance(from, until point) int {
	xDistance := math.Abs(float64(from.x - until.x))
	yDistance := math.Abs(float64(from.y - until.y))
	return int(xDistance + yDistance)
}

func printPath(walls map[point]bool, sizeX, sizeY int, start, end point, path []point) {
	pathSet := make(map[point]bool)
	for _, c := range path {
		pathSet[c] = true
	}
	for y := 0; y < sizeY; y++ {
		for x := 0; x < sizeX; x++ {
			current := point{x, y}
			if current == start {
				fmt.Print("S")
			} else if current == end {
				fmt.Print("E")
			} else if pathSet[current] {
				fmt.Print("X")
			} else if walls[current] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func readInputFile(location string) (map[point]bool, int, int, point, point) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	sizeY := len(lines)
	sizeX := len(lines[0])
	walls := make(map[point]bool)
	var start, end point

	for y, line := range lines {
		for x, c := range line {
			p := point{x, y}
			switch c {
			case '#':
				walls[p] = true
			case 'S':
				start = p
			case 'E':
				end = p
			}
		}
	}

	return walls, sizeX, sizeY, start, end
}
