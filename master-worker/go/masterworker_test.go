package main

import (
	"reflect"
	"testing"
)

func TestTransposeEvenSplit(t *testing.T) {
	input := [][]int{
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
	}
	want := [][]int{
		{1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2},
		{3, 3, 3, 3, 3},
		{4, 4, 4, 4, 4},
		{5, 5, 5, 5, 5},
	}

	got := Transpose(input, 4)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Transpose(input, 4) = %v, want %v", got, want)
	}
}

func TestTransposeMoreWorkersThanRows(t *testing.T) {
	input := [][]int{{2, 4}, {3, 5}}
	want := [][]int{{2, 3}, {4, 5}}

	got := Transpose(input, 3) // more workers requested than rows available

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Transpose(input, 3) = %v, want %v", got, want)
	}
}

func TestSplitRowsDividesEvenly(t *testing.T) {
	input := make([][]int, 10)
	chunks := splitRows(input, 4)

	if len(chunks) != 4 {
		t.Fatalf("got %d chunks, want 4", len(chunks))
	}
	total := 0
	for _, c := range chunks {
		total += len(c)
	}
	if total != len(input) {
		t.Fatalf("chunks cover %d rows, want %d", total, len(input))
	}
}
