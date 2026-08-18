package main

import (
	"fmt"
	"math/rand"
)

func randomMatrix(rows, cols int) [][]int {
	m := make([][]int, rows)
	for r := range m {
		m[r] = make([]int, cols)
		for c := range m[r] {
			m[r][c] = rand.Intn(10)
		}
	}
	return m
}

func printMatrix(m [][]int) {
	for _, row := range m {
		fmt.Println(row)
	}
}

func main() {
	const numWorkers = 4
	input := randomMatrix(10, 20)
	result := Transpose(input, numWorkers)

	fmt.Println("Input:")
	printMatrix(input)
	fmt.Printf("\nTransposed (via %d worker goroutines):\n", numWorkers)
	printMatrix(result)
}
