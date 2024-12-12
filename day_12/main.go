// Decided to experiment with btree
// Also "Less" is kinda like rust 'cmp'. neat.
package main

import (
	"bufio"
	"fmt"
	"github.com/google/btree"
	"os"
)

type Coordinate struct {
	x int
	y int
}

type CoordItem struct {
	Coord Coordinate
	Rune  rune
}

func (a CoordItem) Less(b btree.Item) bool {
	other := b.(CoordItem)
	if a.Coord.y != other.Coord.y {
		return a.Coord.y < other.Coord.y
	}
	return a.Coord.x < other.Coord.x
}

func main() {
	fmt.Println("Example")
	location := "./input/example.txt"
	cd := readInputFile(location)
	one, two := solve(cd)
	fmt.Println("Part One:", one)
	fmt.Println("Part Two:", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	cd = readInputFile(location)
	one, two = solve(cd)
	fmt.Println("Part One:", one)
	fmt.Println("Part Two:", two)
}

func solve(gd *btree.BTree) (int, int) {
	var totalPartOne int
	var totalPartTwo int
	visited := make(map[Coordinate]bool)

	gd.Ascend(func(i btree.Item) bool {
		coordItem := i.(CoordItem)
		current := coordItem.Coord

		if visited[current] {
			return true
		}

		currentRune := coordItem.Rune

		region := getRegionCells(gd, current, currentRune, visited)
		area := len(region)
		perimeter := countPerimeter(region)
		sides := countSides(region, gd, currentRune)

		totalPartOne += area * perimeter
		totalPartTwo += area * sides

		return true
	})

	return totalPartOne, totalPartTwo
}

// DFS
func getRegionCells(gd *btree.BTree, start Coordinate, r rune, visited map[Coordinate]bool) []Coordinate {
	region := []Coordinate{}
	stack := []Coordinate{start}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[current] {
			continue
		}
		visited[current] = true
		region = append(region, current)

		for _, dir := range []Coordinate{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			next := Coordinate{current.x + dir.x, current.y + dir.y}
			if getRune(gd, next) == r && !visited[next] {
				stack = append(stack, next)
			}
		}
	}
	return region
}

func countPerimeter(region []Coordinate) int {
	inRegion := make(map[Coordinate]bool)
	for _, c := range region {
		inRegion[c] = true
	}

	perimeter := 0
	for _, c := range region {
		neighbors := []Coordinate{
			{c.x, c.y - 1},
			{c.x, c.y + 1},
			{c.x - 1, c.y},
			{c.x + 1, c.y},
		}
		for _, n := range neighbors {
			if !inRegion[n] {
				perimeter++
			}
		}
	}
	return perimeter
}

func countSides(region []Coordinate, gd *btree.BTree, currentRune rune) int {
	var corners int

	for _, c := range region {
		corners += getCorners(c, gd, currentRune)
	}

	return corners
}

func getCorners(c Coordinate, gd *btree.BTree, currentRune rune) int {
	neighbors := []Coordinate{
		{c.x, c.y - 1},     // Up
		{c.x, c.y + 1},     // Down
		{c.x - 1, c.y},     // Left
		{c.x + 1, c.y},     // Right
		{c.x - 1, c.y - 1}, // NW
		{c.x + 1, c.y - 1}, // NE
		{c.x - 1, c.y + 1}, // SW
		{c.x + 1, c.y + 1}, // SE
	}

	up := getRune(gd, neighbors[0])
	down := getRune(gd, neighbors[1])
	left := getRune(gd, neighbors[2])
	right := getRune(gd, neighbors[3])

	nw := getRune(gd, neighbors[4])
	ne := getRune(gd, neighbors[5])
	sw := getRune(gd, neighbors[6])
	se := getRune(gd, neighbors[7])

	corners := 0
	// Top left
	if up != currentRune && left != currentRune {
		corners++
	}

	// Top right
	if up != currentRune && right != currentRune {
		corners++
	}

	// Bottom leftt
	if down != currentRune && left != currentRune {
		corners++
	}

	// Bottom right
	if down != currentRune && right != currentRune {
		corners++
	}

	// Bottom left concave
	if sw != currentRune && left == currentRune && down == currentRune {
		corners++
	}

	// Bottom right concave
	if se != currentRune && right == currentRune && down == currentRune {
		corners++
	}

	// Top left concave
	if nw != currentRune && left == currentRune && up == currentRune {
		corners++
	}

	// Top right concave
	if ne != currentRune && right == currentRune && up == currentRune {
		corners++
	}

	return corners
}

func getRune(gd *btree.BTree, coord Coordinate) rune {
	item := gd.Get(CoordItem{Coord: coord, Rune: '.'})
	if item == nil {
		return '.'
	}
	return item.(CoordItem).Rune
}

func readInputFile(location string) *btree.BTree {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	out := btree.New(2)

	var y int
	for scanner.Scan() {
		line := scanner.Text()
		for i, c := range line {
			out.ReplaceOrInsert(CoordItem{Coordinate{i, y}, c})
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(2)
	}

	return out
}
