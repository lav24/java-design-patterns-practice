# Master-Worker — Interview Q&A

Real questions, how to approach each one out loud, and a model answer.
Pair this with `EXECUTION_FLOW.md` for the mechanics.

---

**Q: What is the Master-Worker pattern, and what problem does it solve?**

*How to answer:* One-line definition, then the specific *shape* of problem
it fits — interviewers are often really asking "do you know when NOT to
use this."

> A single, known, fixed batch of work gets split into independent pieces;
> a master hands one piece to each of several workers running in parallel,
> waits for every one of them to finish, then combines their results into
> one final answer. It fits problems that can be decomposed into
> pieces that don't depend on each other's results — matrix operations, batch
> image processing, map-reduce-style aggregation. If piece 2 needed piece
> 1's output first, you couldn't parallelize them this way at all.

---

**Q: How do you guarantee the final combined result comes out correct and
consistent, given that workers finish in an unpredictable order?**

*How to answer:* This is the single most important thing to get right in
this pattern — lead with the rule, then explain why the naive approach
breaks.

> The combining step never looks at "who finished first" — it always
> walks through the original chunks in their *original* order and stitches
> the results back together that way, using each worker's known index or
> id, not the order results happened to arrive in. If you combined by
> arrival order instead, you'd get a scrambled answer that's different
> from run to run — a real, subtle bug, because the program would still
> print *a* result, just occasionally the wrong one, which is much harder
> to catch than a crash.

---

**Q: In the Java version, how is it safe for multiple worker threads to
report results back to the same Master object at the same time?**

*How to answer:* Name the actual mechanism (a thread-safe collection), then
be ready to contrast it with the Go version, which needs no such thing —
that contrast is the more interesting part.

> Results get collected into a `Hashtable` keyed by worker id, and
> `Hashtable`'s methods are synchronized, so concurrent `put` calls from
> different worker threads can't corrupt it. The main thread only reads the
> final combined result after calling `join()` on every worker, which
> creates a happens-before relationship — so there's no risk of reading a
> half-finished result either.

**Follow-up you should expect: "Is there a simpler way to do this without a
thread-safe map?"**

> Yes — if every worker only ever writes to its own reserved slot in a
> pre-sized array or list (by index), there's no shared mutable state
> between workers at all, so you don't need any locking or a
> thread-safe collection. That's actually what the Go version does: each
> goroutine writes to `results[i]`, and since no two goroutines ever touch
> the same index, it's race-free without a mutex — I verified that with
> `go test -race`.

---

**Q: How is this different from Producer-Consumer?**

*How to answer:* Name the one structural difference, not a list of minor
ones — a crisp one-sentence contrast reads as more confident than a long
list.

> Master-Worker has a fixed, known amount of work — someone waits until
> every last piece is done, and then the whole thing is finished, a
> one-shot operation. Producer-Consumer has no concept of "done" at all —
> it just keeps running indefinitely, handling whatever new work shows up,
> for as long as the program is alive.

---

**Q: How is this different from Fork/Join?**

*How to answer:* This tests whether you actually know multiple concurrency
patterns rather than one. Name the structural distinction: recursive
subdivision vs. one flat split.

> Fork/Join is typically recursive — a task splits itself into
> subtasks, which can themselves split further, forming a tree, often with
> work-stealing between idle threads. Master-Worker here is a flat,
> single-level split — the master divides once, hands pieces directly to a
> fixed set of workers, and none of them subdivide further or steal work
> from each other.

---

**Q: This codebase originally had an abstract `Master`/`Worker`/`Input`/
`Result` class hierarchy with exactly one concrete implementation each. Why
did you collapse them into single concrete classes?**

*How to answer:* State the principle, not just "it looked cleaner" — the
principle is what an interviewer is actually checking for.

> An abstract base class only pays for itself if there's more than one
> concrete subclass sharing it. With exactly one subclass, the abstract
> layer added a file and an indirection to read through without adding any
> real flexibility — nothing else in the codebase depended on being able to
> swap in a second `Master` implementation. I merged each pair into one
> concrete class, kept every public method signature the existing tests
> called, and confirmed the full test suite still passed unchanged.

---

**Q: What would you change about this implementation for production use,
versus a teaching example?**

*How to answer:* Name a concrete, specific change, not "make it more
robust" — specificity is what separates a real engineer's answer from a
rehearsed one.

> I'd replace one-thread-per-chunk with a fixed-size thread pool (or a
> goroutine pool in Go), so a job with thousands of chunks doesn't spawn
> thousands of OS threads at once. I'd also add a timeout per worker, so one
> stuck piece of work can't hang the entire batch forever — right now, if a
> single worker thread never finishes, `join()` blocks forever and nothing
> ever completes.

---

**Q: Walk me through what happens, step by step, from calling `getResult()`
to getting the final answer back.**

*How to answer:* This is a "trace the code" question — answer it as a
numbered list, out loud, exactly like you'd trace it on a whiteboard. Don't
skip the barrier step; interviewers listen for whether you mention it.

> `getResult` calls the master's `doWork`, which splits the input into
> row-chunks and hands one to each worker, then starts all worker threads at
> once. Each worker transposes its own chunk with no knowledge of the
> others, then reports its result back to the master, which stores it keyed
> by worker id. Meanwhile, the original thread is sitting in a loop calling
> `join()` on every worker — that's the barrier, it doesn't proceed until
> every single one has finished. Once that wait is over, `getFinalResult`
> hands back the matrix that was already assembled — combined in original
> chunk order — the moment the last worker's result came in.
