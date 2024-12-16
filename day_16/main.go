// Not super fast; Could probably optimize, but for AOC purposes it's fine
package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

func main() {
	fmt.Println("Example")
	exampleLocation := "./input/example.txt"
	maze, start, width, height := readInputFile(exampleLocation)
	one, two := solve(maze, start, width, height)
	fmt.Println("Part One:", one)
	fmt.Println("Part Two:", two)

	fmt.Println("\nInput")
	inputLocation := "./input/input.txt"
	maze, start, width, height = readInputFile(inputLocation)
	one, two = solve(maze, start, width, height)
	fmt.Println("Part One:", one)
	fmt.Println("Part Two:", two)
}

type Coordinate struct {
	x, y int
}

type Maze map[Coordinate]rune

type Path struct {
	current  Coordinate
	previous map[Coordinate]bool
	cost     int
	lastDir  string
	index    int
}

// Decided to try and implement a priority queue with a heap, I'm really Go'ing now
type PriorityQueue []*Path

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	path := x.(*Path)
	*pq = append(*pq, path)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	path := old[n-1]
	*pq = old[0 : n-1]
	return path
}

// Straightforward Dijkstra
func solve(maze Maze, start Coordinate, width, height int) (int, int) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	heap.Push(pq, &Path{current: start, cost: 0, lastDir: "E", previous: map[Coordinate]bool{start: true}})
	bestCost := make(map[Coordinate]map[string]int)
	bestCost[start] = map[string]int{"E": 0}

	directions := map[string]Coordinate{
		"N": {0, -1},
		"S": {0, 1},
		"W": {-1, 0},
		"E": {1, 0},
	}

	lowestCost := math.MaxInt32
	var uniqueCoordinates map[Coordinate]bool = map[Coordinate]bool{start: true}

	for pq.Len() > 0 {
		curPath := heap.Pop(pq).(*Path)
		current := curPath.current

		if curPath.cost > lowestCost {
			continue
		}

		if maze[current] == 'E' {
			if curPath.cost < lowestCost {
				lowestCost = curPath.cost
				// Lower cost = new best path, so reset the unique coordinates
				uniqueCoordinates = make(map[Coordinate]bool)
			}

			if curPath.cost == lowestCost {
				for k := range curPath.previous {
					uniqueCoordinates[k] = true
				}
			}
			continue
		}

		if curPath.cost > bestCost[current][curPath.lastDir] {
			continue
		}

		for dir, offset := range directions {
			next := Coordinate{current.x + offset.x, current.y + offset.y}
			if next.x < 0 || next.x >= width || next.y < 0 || next.y >= height || maze[next] == '#' {
				continue
			}

			stepCost := 1
			if curPath.lastDir != dir {
				stepCost += 1000
			}
			newCost := curPath.cost + stepCost

			if _, exists := bestCost[next]; !exists {
				bestCost[next] = make(map[string]int)
			}

			if cost, seen := bestCost[next][dir]; !seen || newCost <= cost {
				// Keep track of visited nodes, so in part two we can determine the number of unique nodes along our paths
				prev := map[Coordinate]bool{}
				for k, v := range curPath.previous {
					prev[k] = v
				}
				prev[next] = true

				bestCost[next][dir] = newCost
				heap.Push(pq, &Path{current: next, cost: newCost, lastDir: dir, previous: prev})
			}
		}
	}

	printMaze(maze, width, height, uniqueCoordinates)
	return lowestCost, len(uniqueCoordinates)
}

func printMaze(maze Maze, width, height int, partTwo map[Coordinate]bool) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if _, exists := partTwo[Coordinate{x, y}]; exists {
				fmt.Print("O")
			} else {
				fmt.Print(string(maze[Coordinate{x, y}]))
			}
		}
		fmt.Println()
	}
}

func readInputFile(location string) (Maze, Coordinate, int, int) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	maze := Maze{}
	deer := Coordinate{}

	var y int
	for scanner.Scan() {
		line := scanner.Text()
		for i, char := range line {
			if char == 'S' {
				deer = Coordinate{i, y}
				maze[Coordinate{i, y}] = '.'
			} else {
				maze[Coordinate{i, y}] = char
			}
		}
		y++
	}

	return maze, deer, len(maze), y
}
