package main

import "sync"

// eagerSingleton and lazySingleton are unexported — no other package can
// write singleton.eagerSingleton{} directly. This is Go's substitute for
// Java's private constructor: there's no constructor to guard, so the
// guarantee lives at the package boundary instead, enforced by what's
// exported (GetEagerInstance/GetLazyInstance) versus what isn't (the
// struct types themselves).

// --- Eager: a package-level variable ---
//
// Go guarantees every package-level variable finishes initializing before
// main() runs and before any other goroutine can observe the package, so
// this is thread-safe with zero extra code — the same free guarantee
// Java's eager `static final` field gets from the JVM's class-init rules.
type eagerSingleton struct {
	id int
}

var eagerInstance = &eagerSingleton{id: 1}

func GetEagerInstance() *eagerSingleton {
	return eagerInstance
}

// --- Lazy: sync.Once ---
//
// sync.Once IS double-checked locking, already written and already
// race-tested in the standard library: Do runs its function exactly once no
// matter how many goroutines call it concurrently, and every other caller
// blocks until that one run finishes. This replaces Java's hand-rolled
// check-lock-check-again dance and its volatile subtlety entirely.
type lazySingleton struct {
	id int
}

var (
	lazyInstance *lazySingleton
	once         sync.Once
)

func GetLazyInstance() *lazySingleton {
	once.Do(func() {
		lazyInstance = &lazySingleton{id: 1}
	})
	return lazyInstance
}
