package main

import (
	"sync"
	"testing"
)

func TestEagerInstanceIsAlwaysTheSame(t *testing.T) {
	a := GetEagerInstance()
	b := GetEagerInstance()
	if a != b {
		t.Fatalf("got two different instances: %p vs %p", a, b)
	}
}

func TestLazyInstanceIsSameAcrossConcurrentCallers(t *testing.T) {
	const goroutines = 100
	results := make([]*lazySingleton, goroutines)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = GetLazyInstance() // all 100 race here at once
		}(i)
	}
	wg.Wait()

	first := results[0]
	for i, r := range results {
		if r != first {
			t.Fatalf("goroutine %d got a different instance: %p vs %p", i, r, first)
		}
	}
}
