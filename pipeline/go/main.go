package main

import (
	"fmt"
	"unicode"
)

// removeAlphabets is a Handler[string, string]: same shape, just a plain function.
func removeAlphabets(input string) string {
	var result []rune
	for _, c := range input {
		if !unicode.IsLetter(c) {
			result = append(result, c)
		}
	}
	return string(result)
}

func removeDigits(input string) string {
	var result []rune
	for _, c := range input {
		if !unicode.IsDigit(c) {
			result = append(result, c)
		}
	}
	return string(result)
}

func toCharArray(input string) []rune {
	return []rune(input)
}

func main() {
	fmt.Println("Creating pipeline")
	stage1 := NewPipeline(Handler[string, string](removeAlphabets))
	stage2 := AddHandler[string, string, string](stage1, removeDigits)
	stage3 := AddHandler[string, string, []rune](stage2, toCharArray)

	input := "GoYankees123!"
	fmt.Printf("Executing pipeline with input: %s\n", input)
	output := stage3.Execute(input)
	fmt.Printf("Pipeline output: %q\n", string(output))
}
