package main

import "testing"

func TestPipelineChainsStagesInOrder(t *testing.T) {
	stage1 := NewPipeline(Handler[string, string](removeAlphabets))
	stage2 := AddHandler[string, string, string](stage1, removeDigits)
	stage3 := AddHandler[string, string, []rune](stage2, toCharArray)

	got := string(stage3.Execute("#H!E(L&L0O%THE3R#34E!"))
	want := "#!(&%#!"

	if got != want {
		t.Fatalf("Execute(...) = %q, want %q", got, want)
	}
}

func TestRemoveAlphabets(t *testing.T) {
	if got := removeAlphabets("Go123!"); got != "123!" {
		t.Fatalf("removeAlphabets(...) = %q, want %q", got, "123!")
	}
}

func TestRemoveDigits(t *testing.T) {
	if got := removeDigits("Go123!"); got != "Go!" {
		t.Fatalf("removeDigits(...) = %q, want %q", got, "Go!")
	}
}
