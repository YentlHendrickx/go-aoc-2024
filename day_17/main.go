package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Example")
	exampleLocation := "./input/example5.txt"
	regA, regB, regC, instructions := readInputFile(exampleLocation)
	one := solve(regA, regB, regC, instructions)
	fmt.Println("Part One:", one)
	// fmt.Println("Part Two:", solve2(instructions))

	fmt.Println("\nInput")
	inputLocation := "./input/input.txt"
	regA, regB, regC, instructions = readInputFile(inputLocation)
	one = solve(regA, regB, regC, instructions)
	fmt.Println("Part One:", one)
	// fmt.Println("Part Two:", two)
}

func solve(a, b, c int, instructions []Instruction) string {
	var out []int

	var ip int

	regA := a
	regB := b
	regC := c

	for {
		if ip >= len(instructions) {
			break
		}

		fmt.Println()
		fmt.Println("New instruction; ip:", ip)

		instruction := instructions[ip]
		fmt.Println("Opcode:", instruction.opcode)
		fmt.Println("Operand:", instruction.operand)
		fmt.Println("Combo operand:", getCombo(regA, regB, regC, instruction.operand))

		fmt.Println("Before:", regA, regB, regC)

		if instruction.opcode == 3 {
			i := jnz(regA, instruction.operand)
			if i != -1 {
				// ip should be set to i; however we don't increment by two but 1 for each operation
				// Thus we need to set ip to i/2
				ip = int(math.Round(float64(i / 2)))
				continue
			}
		} else {
			operation := getOperation(instruction.opcode)
			resA, resB, resC := operation.apply(regA, regB, regC, instruction.operand)

			if instruction.opcode == 5 {
				out = append(out, resA)
			} else {
				regA = resA
				regB = resB
				regC = resC
			}

		}

		fmt.Println("After:", regA, regB, regC)
		ip++
	}

	var result string
	for _, o := range out {
		result += strconv.Itoa(o) + ","
	}

	if len(result) == 0 {
		return ""
	}

	return result[:len(result)-1]
}

func getCombo(regA, regB, regC, operand int) int {
	if operand < 4 {
		return operand
	}

	if operand == 4 {
		return regA
	}

	if operand == 5 {
		return regB
	}

	if operand == 6 {
		return regC
	}

	return -1
}

type Operation interface {
	apply(regA, regB, regC, operand int) (int, int, int)
}

type Instruction struct {
	opcode  int
	operand int
}

func readInputFile(location string) (int, int, int, []Instruction) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Read first line
	scanner.Scan()
	regA, regB, regC := 0, 0, 0
	fmt.Sscanf(scanner.Text(), "Register A: %d", &regA)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "Register B: %d", &regB)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "Register C: %d", &regC)

	// Read instructions
	scanner.Scan()
	scanner.Scan()
	numbers := strings.Split(scanner.Text(), ",")
	numbers[0] = strings.TrimPrefix(numbers[0], "Program: ")

	instructions := []Instruction{}
	for i := 0; i < len(numbers); i += 2 {
		opcode, _ := strconv.Atoi(numbers[i])
		operand, _ := strconv.Atoi(numbers[i+1])

		instructions = append(instructions, Instruction{opcode, operand})
	}

	return regA, regB, regC, instructions
}

func getOperation(opcode int) Operation {
	switch opcode {
	case 0:
		return adv{}
	case 1:
		return bxl{}
	case 2:
		return bst{}
	// case 3:
	// 	return jnz{}
	case 4:
		return bxc{}
	case 5:
		return out{}
	case 6:
		return bdv{}
	case 7:
		return cdv{}
	}

	return nil
}

type adv struct {
	Operation
}

func (o adv) apply(regA, regB, regC, operand int) (int, int, int) {
	combo := getCombo(regA, regB, regC, operand)
	divisor := int(math.Pow(2, float64(combo)))

	if divisor == 0 {
		return regA, regB, regC
	}

	return regA / divisor, regB, regC
}

type bxl struct {
	Operation
}

func (o bxl) apply(regA, regB, regC, operand int) (int, int, int) {
	return regA, regB ^ operand, regC
}

type bst struct {
	Operation
}

func (o bst) apply(regA, regB, regC, operand int) (int, int, int) {
	combo := getCombo(regA, regB, regC, operand)
	return regA, combo % 8, regC
}

func jnz(regA int, operand int) int {
	if regA == 0 {
		return -1
	}

	return operand
}

type bxc struct {
	Operation
}

func (o bxc) apply(regA, regB, regC, _ int) (int, int, int) {
	return regA, regB ^ regC, regC
}

type out struct {
	Operation
}

func (o out) apply(regA, regB, regC, operand int) (int, int, int) {
	combo := getCombo(regA, regB, regC, operand)

	if combo == 0 {
		return 0, regB, regC
	}

	return combo % 8, regB, regC
}

type bdv struct {
	Operation
}

func (o bdv) apply(regA, regB, regC, operand int) (int, int, int) {
	combo := getCombo(regA, regB, regC, operand)
	divisor := int(math.Pow(2, float64(combo)))

	if divisor == 0 {
		return regA, regB, regC
	}

	return regA, regA / divisor, regC
}

type cdv struct {
	Operation
}

func (o cdv) apply(regA, regB, regC, operand int) (int, int, int) {
	combo := getCombo(regA, regB, regC, operand)
	divisor := int(math.Pow(2, float64(combo)))

	if divisor == 0 {
		return regA, regB, regC
	}

	return regA, regB, regA / divisor
}
