# Singleton — Who Calls What, and Why

This is a study reference, not the project's official docs page (that's `README.md`).
Goal: know exactly which file talks to which file, in what order, in both Java and Go,
well enough to draw it on a whiteboard and write it from memory.

## The big picture, in one paragraph

You want exactly one instance of something to exist for the whole life of
the program, and you want one shared way for any code, anywhere, to reach
it. The hard part isn't the concept — it's making that guarantee hold up
when multiple threads ask for the instance at the exact same time, and
(in Java specifically) when someone tries to sneak around your protection
using reflection or serialization. This module doesn't pick one answer —
it shows 6 different ways to build it, because the interesting material is
comparing the trade-offs, not the idea of "one instance" itself.

---

## Java version

### The files and what each one is

| File | What it is |
|---|---|
| `IvoryTower.java` | Eager — builds the instance immediately when the class loads. |
| `EnumIvoryTower.java` | Uses a Java `enum` with one constant as the instance. |
| `ThreadSafeLazyLoadedIvoryTower.java` | Lazy — builds the instance on first request, guarded by a `synchronized` method. |
| `ThreadSafeDoubleCheckLocking.java` | Lazy, but only locks the first few times, using the "check, lock, check again" trick. |
| `InitializingOnDemandHolderIdiom.java` | Lazy, using a nested class that Java doesn't load until first referenced — no locking at all. |
| `BillPughImplementation.java` | The exact same idiom as the file above, under a different name — two textbooks, one technique. |
| `App.java` | Demonstrates calling `getInstance()` on each variant. |

None of these hold a reference to each other — each file is a
**self-contained, independent answer** to the same question, not
collaborating pieces of one system (unlike the last 3 patterns).

### What actually happens, in order, when you run each one

**Eager (`IvoryTower`)**:
1. The very first time any code touches the `IvoryTower` class (even just calling `getInstance()`), the JVM initializes the class.
2. Initializing the class runs `private static final IvoryTower INSTANCE = new IvoryTower();` right then — building the one instance immediately, before `getInstance()` even executes its first line.
3. Every call to `getInstance()` after that just hands back the already-built `INSTANCE`. No checks, no locking, ever.

**Lazy with a lock (`ThreadSafeLazyLoadedIvoryTower`)**:
1. Nothing is built until the first call to `getInstance()`.
2. That call is `synchronized` — if two threads call it at the same instant, one waits for the other to finish before even checking whether the instance exists.
3. Whichever thread gets in first sees `instance == null`, builds it, and returns it. The other thread then enters, sees it's no longer `null`, and returns the same one.
4. Every future call still goes through that same lock — even though nothing needs building anymore, every caller still pays a small locking cost forever.

**Double-checked locking (`ThreadSafeDoubleCheckLocking`)**:
1. First check of `instance == null` happens with **no locking at all** — cheap.
2. If it's already built (the common case, after the first call), return immediately — this is why this version is faster than the plain-synchronized one long-term.
3. If it looks like it might not be built yet, *now* acquire a lock, and check **again** — because another thread might have finished building it in the gap between the first check and acquiring the lock.
4. Only if it's still `null` after the second check does this thread actually build it.

**Holder idiom (`InitializingOnDemandHolderIdiom` / `BillPughImplementation`)**:
1. Calling `getInstance()` for the first time references the nested holder class for the first time.
2. Referencing a class for the first time is exactly the trigger that makes Java initialize it — so the holder class's own instance field gets built right then, guaranteed safe and exactly-once by the same JVM class-loading rules the eager version relies on.
3. No explicit lock or `volatile` needed anywhere — the laziness comes from the holder class simply not being loaded until it's asked for.

### Method-by-method: who calls what

| File | Method | Called by | It calls |
|---|---|---|---|
| `IvoryTower.java` | `getInstance()` | any caller (e.g. `App`, tests) | nothing — returns the field built at class-load time |
| `ThreadSafeLazyLoadedIvoryTower.java` | `getInstance()` | any caller | `new ThreadSafeLazyLoadedIvoryTower()`, only the first time |
| `ThreadSafeDoubleCheckLocking.java` | `getInstance()` | any caller | `new ThreadSafeDoubleCheckLocking()`, only the first time, inside the inner locked check |
| `InitializingOnDemandHolderIdiom.java` | `getInstance()` | any caller | `HelperHolder.INSTANCE` — referencing this triggers `HelperHolder`'s one-time class initialization |
| `BillPughImplementation.java` | `getInstance()` | any caller | `InstanceHolder.instance`, same mechanism as above |
| `EnumIvoryTower.java` | (none needed) | any caller reads `EnumIvoryTower.INSTANCE` directly | nothing |

---

## Go version

### The files and what each one is

| File | What it is |
|---|---|
| `singleton.go` | Both singleton variants: `eagerSingleton` (package variable) and `lazySingleton` (`sync.Once`). |
| `main.go` | Demonstrates both, including 100 goroutines racing for the lazy one. |

Go has no equivalent of "private constructor," because Go structs don't
have constructors at all — `GetEagerInstance`/`GetLazyInstance` are just
regular functions. The guarantee here rests entirely on `eagerSingleton`
and `lazySingleton` being **unexported** (lowercase) — no code outside this
file can write `singleton.eagerSingleton{}` directly, because the package
boundary is Go's access-control mechanism, not a per-type keyword.

### What actually happens, in order, when you run it

**Eager**:
1. Before `main()` runs at all, Go initializes every package-level variable — including `var eagerInstance = &eagerSingleton{id: 1}`.
2. By the time any code can call `GetEagerInstance()`, the instance already exists. The function just returns the pointer.

**Lazy**:
1. `lazyInstance` starts as `nil`; nothing is built yet.
2. The first goroutine to call `GetLazyInstance()` triggers `once.Do(...)`, which runs the build function.
3. Every *other* goroutine that calls `GetLazyInstance()` at the same time — even mid-build — blocks inside `once.Do` until that one build finishes, then all of them (including the one that built it) return the same pointer.
4. Every future call sees `sync.Once` already "used up" and returns instantly, no build, no blocking.

### Method-by-method: who calls what

| File | Function | Called by | It calls |
|---|---|---|---|
| `singleton.go` | `GetEagerInstance()` | `main()` (and any other package code) | nothing — returns the package variable |
| `singleton.go` | `GetLazyInstance()` | `main()`, potentially by many goroutines at once | `once.Do(...)`, which runs the build function at most one time total |
| `main.go` | `main()` | program entry point | both getters; launches 100 goroutines to prove the lazy one is race-free |

---

## Side-by-side: same idea, different word

| What it does | Java | Go |
|---|---|---|
| Preventing outside code from building its own copy | `private` constructor (+ a manual guard, since reflection can bypass `private` alone) | keeping the struct type unexported — no constructor exists to guard in the first place |
| Eager, thread-safe by default | `static final` field, built at class-load | package-level variable, built before `main()` |
| Lazy, thread-safe, hand-rolled | double-checked locking + `volatile` | not needed — `sync.Once` replaces this entirely |
| Lazy, thread-safe, zero locking | the holder idiom (nested class loaded on first reference) | no equivalent — Go doesn't defer package-variable initialization the way Java defers nested-class loading |
| Bulletproof against reflection/deserialization | `enum` singleton | no equivalent — Go has no enum-like construct the runtime specially protects |

---

## Pseudocode you can write from scratch in an interview

```
# Eager
GLOBAL instance = build()          # built once, at startup, before anyone can ask
function getInstance():
    return instance

# Lazy, naive (works, but slow forever after the first call)
GLOBAL instance = null
function getInstance():
    lock()
    if instance is null:
        instance = build()
    unlock()
    return instance

# Lazy, double-checked (fast after the first call)
GLOBAL instance = null
function getInstance():
    if instance is not null:          # cheap check, no lock, handles 99% of calls
        return instance
    lock()
    if instance is null:              # check AGAIN — someone else may have just built it
        instance = build()
    unlock()
    return instance
```

The one line worth memorizing: **"double-checked" means checking before
locking (fast path, no contention) AND checking again after locking (to
avoid building it twice if you lost a race to acquire the lock)."**

---

## Likely interview questions, answered short

See `QA.md` in this same folder for the full interview-style question set
with model answers. Short version of the highlights:

**Why does double-checked locking need `volatile`?**
Without it, another thread could see a reference to an object whose
constructor hasn't actually finished running yet, due to reordering — you'd
be handed a half-built object. `volatile` prevents that reordering.

**Is Singleton considered an anti-pattern?**
Often, yes — it's global mutable state and makes swapping in a test double
harder. Most modern code prefers a dependency-injection container managing
a single instance instead of hand-rolling this.

**How is Go's approach fundamentally different in what it can guarantee?**
Java's guarantee (once you add the guard code) survives an attacker
deliberately trying to break it via reflection. Go's guarantee rests
entirely on "we didn't export the type" — there's no reflection-blocking or
constructor-guarding mechanism to add even if you wanted one.
