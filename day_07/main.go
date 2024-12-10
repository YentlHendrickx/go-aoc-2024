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
	var location string = "./input/example.txt"
	equations := readInputFile(location)

	one, two := solve(equations)
	fmt.Println("PartOne: ", one)
	fmt.Println("PartTwo: ", two)

	fmt.Println("\nInput")
	location = "./input/input.txt"
	equations = readInputFile(location)

	one, two = solve(equations)
	fmt.Println("PartOne: ", one)
	fmt.Println("PartTwo: ", two)
}

func solve(equations []Equation) (int, int) {
	var totalOne int = 0
	var totalTwo int = 0

	for _, eq := range equations {
		// In most cases the first find will actually return the correct result. Optimization is that we can skip the concatenate operation unless we have to
		// So this does result in parsing the equation twice, but in the general case it will be faster since matching_1 is more than twice the count of matching_2
		// It's only like 50ms different from what i can see
		if findOperations(&eq, 1, eq.numbers[0], []func(a, b int) int{add, multiply}) {
			totalOne += eq.result
			totalTwo += eq.result
		} else if findOperations(&eq, 1, eq.numbers[0], []func(a, b int) int{add, multiply, concatenate}) {
			totalTwo += eq.result
		}
	}

	return totalOne, totalTwo
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func concatenate(a, b int) int {
	temp := strconv.FormatInt(int64(a), 10)
	temp += strconv.FormatInt(int64(b), 10)
	res, _ := strconv.Atoi(temp)
	return res
}

func findOperations(equation *Equation, depth int, current int, ops []func(a, b int) int) bool {
	if current > equation.result {
		return false
	} else if depth == len(equation.numbers) {
		return current == equation.result
	} else {
		for _, op := range ops {
			if findOperations(equation, depth+1, op(current, equation.numbers[depth]), ops) {
				return true
			}
		}

		return false
	}
}

type Equation struct {
	result  int
	numbers []int
}

func readInputFile(location string) []Equation {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var out []Equation = make([]Equation, 0)

	for scanner.Scan() {
		line := scanner.Text()
		var eq Equation

		numbers := strings.Fields(line)
		for _, num := range numbers {
			if eq.result == 0 {
				eq.result, _ = strconv.Atoi(strings.Replace(num, ":", "", -1))
				continue
			}

			n, _ := strconv.Atoi(num)
			eq.numbers = append(eq.numbers, n)
		}

		out = append(out, eq)
	}

	return out
}
