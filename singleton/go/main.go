package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Eager singleton:")
	a := GetEagerInstance()
	b := GetEagerInstance()
	fmt.Printf("  two calls return the same pointer? %v\n", a == b)

	fmt.Println("Lazy singleton, requested by 100 concurrent goroutines:")
	const goroutines = 100
	results := make([]*lazySingleton, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = GetLazyInstance()
		}(i)
	}
	wg.Wait()

	allSame := true
	for _, r := range results {
		if r != results[0] {
			allSame = false
			break
		}
	}
	fmt.Printf("  all 100 goroutines got the same pointer? %v\n", allSame)
}
