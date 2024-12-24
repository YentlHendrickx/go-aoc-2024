package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Gate struct {
	in1, in2 string
	op       string
	out      string
}

type States map[string]int

func main() {
	fmt.Println("Example")
	solvePart("./input/example.txt", "")

	fmt.Println("\nExample 2")
	solvePart("./input/example_2.txt", "")

	fmt.Println("\nInput")
	// Yeah so doing this with code is hard af
	// I just graphvized this!
	solvePart("./input/input.txt", "./input/input.dot")
}

func solvePart(location, dotOutput string) {
	part1 := solve(location, dotOutput, false)
	fmt.Println("Part 1:", part1)

	// Part 2
	solve(location, "./input/input_fixed.dot", true)
}

func solve(location, dotOutput string, partTwo bool) int {
	states, gates := readInputFile(location)

	if partTwo {
		// After GraphViz this came out!
		// Swap z10 <-> ggn
		// Swap z39 <-> twr
		// Swap z32 <-> grm
		// Swap ndw <-> jcb
		swapPairs := []string{"z10", "ggn", "z39", "twr", "z32", "grm", "ndw", "jcb"}
		for i := 0; i < len(swapPairs); i += 2 {
			swapGateOutputs(swapPairs[i], swapPairs[i+1], gates)
		}

		// Sort swappairs by alphabetical order and printer comma separated
		sort.Strings(swapPairs)
		fmt.Println(strings.Join(swapPairs, ","))
	}

	expected, actual := solveStructure(states, gates)

	if partTwo {
		if expected == actual {
			fmt.Println("Expected and actual values are the same")
		} else {
			fmt.Printf("Expected: \n%b\n%b\n", expected, actual)
		}
	}

	if err := writeGraphviz(dotOutput, gates); err != nil {
		fmt.Printf("Error writing Graphviz file: %v\n", err)
	}

	return actual
}

func swapGateOutputs(a, b string, gates []Gate) {
	for i, g := range gates {
		if g.out == a {
			gates[i].out = b
		} else if g.out == b {
			gates[i].out = a
		}
	}
}

func solveStructure(states States, gates []Gate) (int, int) {
	for _, gate := range gates {
		ensureWireExists(states, gate.in1)
		ensureWireExists(states, gate.in2)
		ensureWireExists(states, gate.out)
	}

	for {
		anyUnsolved := false
		for _, g := range gates {
			v1 := states[g.in1]
			v2 := states[g.in2]
			if v1 == -1 || v2 == -1 {
				anyUnsolved = true
				continue
			}

			switch g.op {
			case "AND":
				states[g.out] = v1 & v2
			case "OR":
				states[g.out] = v1 | v2
			case "XOR":
				states[g.out] = v1 ^ v2
			}
		}
		if !anyUnsolved {
			break
		}
	}

	var xVal, yVal, zVal int
	for wireName, val := range states {
		if val < 0 {
			continue
		}
		idx, _ := strconv.Atoi(wireName[1:])
		switch wireName[0] {
		case 'x':
			xVal += val << idx
		case 'y':
			yVal += val << idx
		case 'z':
			zVal += val << idx
		}
	}
	return zVal, xVal + yVal
}

func ensureWireExists(states States, wireName string) {
	if wireName == "" {
		return
	}
	if _, ok := states[wireName]; !ok {
		states[wireName] = -1
	}
}

func readInputFile(location string) (States, []Gate) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var gates []Gate
	states := make(States)

	readingStates := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			readingStates = false
			continue
		}

		if readingStates {
			state, value := getState(line)
			states[state] = value
		} else {
			gate := parseGate(line)
			gates = append(gates, gate)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(2)
	}

	return states, gates
}

func parseGate(line string) Gate {
	parts := strings.Split(line, " -> ")
	out := parts[1]

	left := strings.Split(parts[0], " ")
	return Gate{in1: left[0], op: left[1], in2: left[2], out: out}
}

func getState(line string) (string, int) {
	parts := strings.Split(line, ":")
	wireName := strings.TrimSpace(parts[0])
	val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return wireName, val
}

func writeGraphviz(filename string, gates []Gate) error {
	if filename == "" {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "digraph G {")
	fmt.Fprintln(f, "  rankdir=\"LR\";") // left-to-right layout

	for i, g := range gates {
		gateName := fmt.Sprintf("G%d", i)
		label := g.op
		if label == "" {
			label = "ASSIGN"
		}

		fmt.Fprintf(f, "  %s [shape=box, label=\"%s\"];\n", gateName, label)

		if g.in1 != "" {
			fmt.Fprintf(f, "  \"%s\" -> %s;\n", g.in1, gateName)
		}

		if g.in2 != "" {
			fmt.Fprintf(f, "  \"%s\" -> %s;\n", g.in2, gateName)
		}
		fmt.Fprintf(f, "  %s -> \"%s\";\n", gateName, g.out)
	}

	fmt.Fprintln(f, "}")
	return nil
}
