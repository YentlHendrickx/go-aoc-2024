package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	fmt.Println("Example")
	solvePart("./input/example.txt")

	fmt.Println("\nInput")
	solvePart("./input/input.txt")
}

func solvePart(location string) {
	one, two := solve(location)
	fmt.Println("Part 1:", one)
	fmt.Println("Part 2:", two)
}

type N3Pair struct {
	n1 string
	n2 string
	n3 string
}

func solve(location string) (int, string) {
	links := readInputFile(location)

	adj := makeAdjacency(links)

	allNodes := make([]string, 0, len(adj))
	for node := range adj {
		allNodes = append(allNodes, node)
	}
	sort.Strings(allNodes)

	// Bron–Kerbosch wants sets; we'll use slices
	r := []string{}                      // Current clique being built
	p := append([]string{}, allNodes...) // Potential (candidate) nodes
	x := []string{}                      // Excluded nodes

	var allCliques [][]string
	bronKerbosch(r, p, x, adj, &allCliques)

	var largest []string
	var p1 int
	uniquePairs := make(map[string]N3Pair)
	for _, clique := range allCliques {
		if len(clique) > len(largest) {
			largest = clique
		}

		if len(clique) >= 3 {
			p1 += partOne(clique, uniquePairs)
		}
	}

	password := strings.Join(largest, ",")
	return p1, password
}

func partOne(set []string, unique map[string]N3Pair) int {
	// If the set is longer than 3, we need to find all possible combinations
	pairs := make([]N3Pair, 0)
	if len(set) > 3 {
		// Loop unrolling will save me here; I don't want to write a recursive function
		for i := 0; i < len(set); i++ {
			for j := i + 1; j < len(set); j++ {
				for k := j + 1; k < len(set); k++ {
					key := createKey(set[i], set[j], set[k])
					if _, ok := unique[key]; !ok {
						unique[key] = N3Pair{set[i], set[j], set[k]}
						pairs = append(pairs, N3Pair{set[i], set[j], set[k]})
					}
				}
			}
		}
	} else {
		if _, ok := unique[createKey(set[0], set[1], set[2])]; !ok {
			unique[createKey(set[0], set[1], set[2])] = N3Pair{set[0], set[1], set[2]}
			pairs = append(pairs, N3Pair{set[0], set[1], set[2]})
		}
	}

	// Next we check if pairs have a 't' at the start of any of the computers
	res := 0
	for _, pair := range pairs {
		if pair.n1[0] == 't' || pair.n2[0] == 't' || pair.n3[0] == 't' {
			res++
		}
	}

	return res
}

func createKey(n1, n2, n3 string) string {
	alphabetized := []string{n1, n2, n3}
	sort.Strings(alphabetized)
	return strings.Join(alphabetized, ",")
}

// Intuition told me there was an algorithm for this; didn't dissapoint
// I should've probably tried to implement it myself before looking up the algo, as it's quite easy
// But I'm glad I found it, as it's a nice algo to know
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm
// Comments are not all mine
func bronKerbosch(r, p, x []string, adj map[string]map[string]bool, allCliques *[][]string) {
	// If no more candidates (P) and no more nodes to exclude (X),
	// then R is a *maximal* clique. We record it.
	if len(p) == 0 && len(x) == 0 {
		clique := make([]string, len(r))
		copy(clique, r)
		*allCliques = append(*allCliques, clique)
		return
	}

	pSnapshot := append([]string{}, p...)

	for _, v := range pSnapshot {
		// N(v) = all neighbors of v
		neighborsOfV := neighbors(v, adj)

		// Intersection of P with neighbors(v)
		pCapN := intersect(p, neighborsOfV)

		// Intersection of X with neighbors(v)
		xCapN := intersect(x, neighborsOfV)

		// Recurse with R ∪ {v}, P ∩ N(v), X ∩ N(v)
		bronKerbosch(append(r, v), pCapN, xCapN, adj, allCliques)

		// Move v from P to X
		p = removeFromSlice(p, v)
		x = append(x, v)
	}
}

func makeAdjacency(links Link) map[string]map[string]bool {
	adj := make(map[string]map[string]bool)
	for node, neighbors := range links {
		if adj[node] == nil {
			adj[node] = make(map[string]bool)
		}
		for _, n := range neighbors {
			adj[node][n] = true
		}
	}
	return adj
}

func neighbors(node string, adj map[string]map[string]bool) []string {
	nmap := adj[node]
	result := make([]string, 0, len(nmap))
	for neigh := range nmap {
		result = append(result, neigh)
	}
	return result
}

func intersect(a, b []string) []string {
	setB := make(map[string]bool, len(b))
	for _, val := range b {
		setB[val] = true
	}
	var out []string
	for _, val := range a {
		if setB[val] {
			out = append(out, val)
		}
	}
	return out
}

func removeFromSlice(s []string, elem string) []string {
	var out []string
	for _, v := range s {
		if v != elem {
			out = append(out, v)
		}
	}
	return out
}

type Link map[string][]string

func readInputFile(location string) Link {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	res := make(Link)

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.Split(line, "-")
		left := split[0]
		right := split[1]

		if _, ok := res[left]; !ok {
			res[left] = []string{}
		}

		res[left] = append(res[left], right)

		// Do the same for right
		if _, ok := res[right]; !ok {
			res[right] = []string{}
		}

		res[right] = append(res[right], left)
	}

	return res
}
