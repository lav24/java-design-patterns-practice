# Master-Worker — Who Calls What, and Why

This is a study reference, not the project's official docs page (that's `README.md`).
Goal: know exactly which file talks to which file, in what order, in both Java and Go,
well enough to draw it on a whiteboard and write it from memory.

## The big picture, in one paragraph

You have one big job (in this project: transposing a matrix). One "boss"
splits the job into a handful of smaller pieces, hands one piece to each of
several "helpers" who all work on their piece at the same time, then waits
until every single helper has finished before combining all their answers
into the one final answer. Unlike Producer-Consumer, this isn't a
never-ending stream — it's one job, split up, done, combined, finished.

---

## Java version

### The files and what each one is

| File | What it is |
|---|---|
| `ArrayInput.java` | Wraps the input matrix. Knows how to cut itself into row-chunks, one chunk per helper. |
| `ArrayResult.java` | Wraps a matrix — either one helper's partial answer, or the final combined answer. Just data. |
| `ArrayUtilityMethods.java` | Small helpers: make a random matrix, print a matrix, compare two matrices (used by tests and by `App`). |
| `system/ArrayTransposeMasterWorker.java` | The front door. This is the only class the outside world (`App`) talks to. |
| `system/systemmaster/ArrayTransposeMaster.java` | The boss. Splits the work, hands out chunks, waits for everyone, combines the answers. |
| `system/systemworkers/ArrayTransposeWorker.java` | A helper. Runs on its own thread, transposes the one chunk it was given, reports back. |
| `App.java` | The program's starting point. |

### Who holds a reference to whom (the relationships)

```
App
 └── creates ONE  ArrayTransposeMasterWorker
                     └── creates ONE  ArrayTransposeMaster
                                        └── creates FOUR ArrayTransposeWorker objects
                                                          (each one holds a reference back to the Master,
                                                           so it knows who to report its answer to)

ArrayTransposeMasterWorker  ---holds a reference to--->  ArrayTransposeMaster
ArrayTransposeMaster        ---holds a reference to--->  all 4 ArrayTransposeWorkers
ArrayTransposeWorker        ---holds a reference back to--->  ArrayTransposeMaster

(Workers never talk to each other. Every arrow above either points at the
 Master, or comes from the Master. The Master is the only thing every
 other piece knows about.)
```

### What actually happens, in order, when you run it

1. `App.main()` creates one `ArrayTransposeMasterWorker`. That, in turn, immediately creates one `ArrayTransposeMaster`, which immediately creates its 4 `ArrayTransposeWorker` objects — but none of them are running yet, they're just sitting there waiting to be given work.
2. `App.main()` builds a random 10×20 matrix and wraps it in an `ArrayInput`.
3. `App.main()` calls `masterWorker.getResult(input)`. This is the one call that kicks off everything else — and `App` doesn't do anything else itself; it just waits for this call to return.
4. Inside that call, the Master calls `input.divideData(4)`. This slices the 10 rows into 4 row-chunks (3 rows, 3 rows, 2 rows, 2 rows — as even as possible).
5. The Master hands chunk 1 to worker 1, chunk 2 to worker 2, and so on, then calls `.start()` on each worker. This is the moment 4 threads actually begin running, all at once, independently.
6. Each worker, on its own thread, transposes only the chunk it was given — it has no idea the other 3 workers or chunks even exist. When it's done, it hands its answer back to the Master.
7. As soon as a worker hands back its answer, the Master checks: "have all 4 answers come in yet?" If yes (which will only be true for whichever worker happens to finish last), it combines them into the final matrix right then and there.
8. Meanwhile, back on the original thread, the Master is sitting in a loop waiting for all 4 threads to finish (this is the "wait for everyone" step). Once that wait is over, it hands back the finished, combined matrix.
9. `App.main()` gets the result back and prints both the original and transposed matrices.

### Method-by-method: who calls what

| File | Method | Called by | It calls |
|---|---|---|---|
| `ArrayInput.java` | `divideData(num)` | `ArrayTransposeMaster.doWork()` | its own helper `makeDivisions(...)` |
| `ArrayTransposeMasterWorker.java` | `getResult(input)` | `App.main()` | `master.doWork(input)`, then `master.getFinalResult()` |
| `ArrayTransposeMaster.java` | `doWork(input)` | `ArrayTransposeMasterWorker.getResult()` | `input.divideData(4)`, each `worker.setReceivedData(...)`, each `worker.start()`, each `worker.join()` |
| `ArrayTransposeMaster.java` | `receiveData(data, worker)` | `ArrayTransposeWorker.run()` — called from a worker's own thread, not the main thread | `aggregateData()`, but only once the last answer has arrived |
| `ArrayTransposeMaster.java` | `aggregateData()` | itself, from inside `receiveData(...)` | nothing external — just combines the stored answers |
| `ArrayTransposeWorker.java` | `run()` | the Java runtime, automatically, the moment `.start()` is called | `executeOperation()`, then `master.receiveData(...)` |
| `ArrayTransposeWorker.java` | `executeOperation()` | its own `run()` | nothing — pure calculation on the chunk it was given |
| `App.java` | `main()` | the program starts here | creates `ArrayTransposeMasterWorker`, builds input, calls `getResult(...)`, prints the result |

A detail worth remembering for interviews: **step 7 runs on a worker's
thread, not the main thread** — whichever worker happens to finish last is
the one that actually does the combining. The main thread only finds out
about the finished result afterward, once its wait (`join`) is over.

---

## Go version

### The files and what each one is

| File | What it is |
|---|---|
| `masterworker.go` | Everything: the boss function (`Transpose`), the splitting function (`splitRows`), the per-chunk work (`transposeChunk`), and the combining function (`mergeColumns`). |
| `main.go` | The program's starting point. |

Go doesn't need separate `Master`/`Worker` classes — a "worker" here is just
a small piece of code launched with `go`, and "waiting for everyone to
finish" is one built-in tool (`sync.WaitGroup`) instead of a hand-written
loop calling `.join()` on each one individually.

### Who calls whom

```
main()
 └── calls  Transpose(input, 4)
               ├── calls  splitRows(input, 4)              — cuts the matrix into up to 4 chunks
               ├── launches ONE goroutine per chunk         — each runs transposeChunk on its own chunk
               ├── waits for all of them (WaitGroup)
               └── calls  mergeColumns(results)             — combines every chunk's answer
```

### What actually happens, in order, when you run it

1. `main()` builds a random 10×20 matrix and calls `Transpose(input, 4)`. Everything else happens inside that one call.
2. `Transpose` calls `splitRows`, which cuts the 10 rows into (up to) 4 contiguous chunks.
3. For each chunk, `Transpose` launches one goroutine — a small independent unit of work — that transposes just that chunk and writes its answer into its own reserved slot in a shared list. Each goroutine only ever touches its own slot, so there's no risk of two of them colliding.
4. `Transpose` then waits (`wg.Wait()`) until every goroutine has finished. This is the same idea as Java's "wait for everyone", just done with one built-in tool instead of a loop calling `.join()` four times.
5. Once every goroutine is done, `Transpose` calls `mergeColumns`, which stitches all 4 partial answers together, in the same left-to-right order the original chunks were in — not in whatever order the goroutines happened to finish.
6. `Transpose` returns the finished matrix to `main()`, which prints both the original and transposed matrices.

### Method-by-method: who calls what

| File | Function | Called by | It calls |
|---|---|---|---|
| `masterworker.go` | `Transpose(input, numWorkers)` | `main()` | `splitRows(...)`, launches goroutines running `transposeChunk(...)`, `mergeColumns(...)` |
| `masterworker.go` | `splitRows(input, numWorkers)` | `Transpose(...)` | nothing — pure slicing logic |
| `masterworker.go` | `transposeChunk(rows)` | one goroutine, launched inside `Transpose(...)` | nothing — pure calculation on the one chunk it was given |
| `masterworker.go` | `mergeColumns(results)` | `Transpose(...)`, after every goroutine has finished | nothing — pure combining logic |
| `main.go` | `main()` | the program starts here | `randomMatrix(...)`, `Transpose(...)`, `printMatrix(...)` |

---

## Side-by-side: same idea, different word

| What it does | Java | Go |
|---|---|---|
| The boss that splits and combines | `ArrayTransposeMaster` class | the `Transpose` function |
| A helper doing one chunk of work | `ArrayTransposeWorker`, running on its own `Thread` | one goroutine, started with `go` |
| Starting a helper | `worker.start()` | `go func(...) {...}(...)` |
| Waiting for every helper to finish | a loop calling `worker.join()` for each one | one call: `wg.Wait()` |
| Collecting each helper's answer safely | a synchronized map (`Hashtable`) keyed by worker id | each goroutine writes to its own reserved slot in a plain list — no locking needed |
| Combining the answers in the right order | loops over the Master's own worker list (by id), not by arrival order | loops over the results list by index (chunk order), not by finish order |

---

## Pseudocode you can write from scratch in an interview

```
function splitIntoChunks(input, numWorkers):
    divide input's rows into numWorkers roughly-equal pieces
    return the list of pieces

function doOnePieceOfWork(chunk):
    # this is the ONLY thing a worker knows how to do — nothing about the
    # other workers, nothing about the original full input
    return transform(chunk)

function master(input, numWorkers):
    chunks = splitIntoChunks(input, numWorkers)
    results = empty list, one reserved slot per chunk

    start one independent worker per chunk:
        worker i: results[i] = doOnePieceOfWork(chunks[i])

    wait until every worker above has finished        # the barrier

    return combine(results)   # combine IN CHUNK ORDER, not finish order
```

The one line worth memorizing for an interview: **"reserve a slot per
worker up front, let each worker write only to its own slot, and combine by
that original index — never by whichever order they happened to finish
in."** That single idea is what makes the output deterministic even though
the workers themselves finish in an unpredictable order.

---

## Likely interview questions, answered short

**What problem does this solve?**
A single big job that can be broken into independent pieces, where doing
the pieces at the same time instead of one-after-another finishes the whole
job faster.

**Why does the final answer always come out the same, even though the
workers finish in a different order every time you run it?**
Because the combining step never looks at "who finished first" — it always
walks through the chunks in their *original* order and stitches them back
together that way. If it combined in arrival order instead, the answer
would come out scrambled differently on different runs — a real, easy bug
to introduce and hard to notice, since it would still print *a* matrix,
just the wrong one about half the time.

**How is this different from Producer-Consumer?**
Master-Worker has one known, fixed batch of work — someone waits until
every last piece is finished, then the whole thing is over. Producer-
Consumer never "finishes" on its own — it just keeps running, taking in
whatever new work shows up, for as long as the program is alive.

**What has to be true about the pieces of work for this pattern to make
sense?**
They have to be independent of each other — no piece needs to know the
result of another piece to do its own job. If piece 2 needed piece 1's
answer first, you couldn't run them at the same time, and this pattern
wouldn't fit.

**Why does the Java version need a separate synchronized map, but Go
doesn't need any locking at all?**
In the Java version, every worker calls the *same* method
(`receiveData`) on the *same* shared object at possibly the same instant,
so that method needs its own protection against being run by two threads
at once — hence a thread-safe map. In the Go version, each worker only
ever writes to its own reserved slot in the results list; two goroutines
never touch the same memory, so there's nothing to protect.

**What would you change for production use instead of a learning example?**
Use a thread pool / goroutine pool with a fixed size instead of one
thread per chunk (so a job with 10,000 chunks doesn't start 10,000
threads), and add a timeout per worker so one stuck piece of work can't
hang the entire job forever.
