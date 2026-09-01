# Pipeline — Who Calls What, and Why

This is a study reference, not the project's official docs page (that's `README.md`).
Goal: know exactly which file talks to which file, in what order, in both Java and Go,
well enough to draw it on a whiteboard and write it from memory.

## The big picture, in one paragraph

You have one value that needs to go through several small transformations,
one after another — like an assembly line. Each station only knows how to
do its own one job: take something in, hand something out. You snap the
stations together in order, and the whole line becomes one thing you can
run by feeding it the very first input. Nothing here runs at the same time
as anything else — it's ordered, step 1 then step 2 then step 3, not
fanned out. That's the difference from Master-Worker and Producer-Consumer:
this pattern isn't about concurrency at all, it's about composing a chain
of transformations in a type-safe way.

---

## Java version

### The files and what each one is

| File | What it is |
|---|---|
| `Handler.java` | The contract every stage follows: "give me one thing, I'll hand back one thing." |
| `Pipeline.java` | The chain-builder. Snaps stages together and runs the whole chain. |
| `App.java` | Defines the actual stages (as small functions) and wires them into a pipeline. |

Earlier versions of this module had 3 more files — one whole class per
stage (`RemoveAlphabetsHandler`, `RemoveDigitsHandler`,
`ConvertToCharArrayHandler`). Since `Handler` has exactly one method,
each of those classes was just one function wearing a class costume — they
now live as small functions directly inside `App.java`.

### Who calls whom (the relationships)

```
App.main()
 ├── builds Handler #1  (removeAlphabets)
 ├── builds Handler #2  (removeDigits)
 ├── builds Handler #3  (toCharArray)
 └── creates ONE Pipeline, snapping all 3 together in order

Pipeline  ---holds---> one Handler
                         (which may itself be built out of several
                          Handlers nested inside each other — see below)
```

### What actually happens, in order, when you run it

1. `App.main()` asks for handler #1 (`removeAlphabets`) and wraps it in a brand-new `Pipeline`. At this point the pipeline can only do one thing: remove alphabet characters.
2. `App.main()` calls `.addHandler(removeDigits)` on that pipeline. This doesn't run anything yet — it returns a **new** pipeline whose one internal handler is "run the old pipeline, then run `removeDigits` on the result."
3. `App.main()` calls `.addHandler(toCharArray)` on *that*. Same thing again: a newer pipeline whose internal handler is "run the previous 2-stage pipeline, then run `toCharArray` on the result."
4. Nothing has actually processed any data yet — building the pipeline is just nesting functions inside each other, three layers deep now.
5. `App.main()` calls `.execute("GoYankees123!")`. This is the one moment any real work happens. It runs the outermost layer, which runs the layer inside it, which runs the layer inside *that* — so all 3 stages fire in order, in one single call: remove letters → `"123!"`, remove digits → `"!"`, convert to char array → `['!']`.
6. The final result comes back out of that one `execute` call and gets printed.

### Method-by-method: who calls what

| File | Method | Called by | It calls |
|---|---|---|---|
| `App.java` | `removeAlphabets()` | `App.main()` (and `PipelineTest`, directly) | nothing — returns a lambda, doesn't run it |
| `App.java` | `removeDigits()` | `App.main()` (and `PipelineTest`) | nothing — returns a lambda |
| `App.java` | `toCharArray()` | `App.main()` (and `PipelineTest`) | nothing — returns a method reference |
| `Pipeline.java` | `addHandler(next)` | `App.main()`, chained fluently | builds a new lambda that calls the old handler, then `next` |
| `Pipeline.java` | `execute(input)` | `App.main()` | the pipeline's one (possibly deeply-nested) handler |
| `App.java` | `main()` | program entry point | the 3 stage-builder methods, `Pipeline`'s constructor, `addHandler`, `execute` |

**The one detail worth remembering**: `execute` always looks like it's
calling just one handler, no matter how many stages you've chained. All the
real complexity is hidden inside the nested lambda that `addHandler` built
earlier — `execute` itself never changes.

---

## Go version

### The files and what each one is

| File | What it is |
|---|---|
| `pipeline.go` | Defines `Handler` (just a function type, not an interface — Go doesn't need one) and `Pipeline` (the chain-builder). |
| `main.go` | Defines the actual stage functions and wires them into a pipeline. |

Go doesn't need a `Handler` *interface* at all, because Go functions are
already values you can pass around. `Handler[I, O]` here is just a **name**
for "a function that takes an `I` and returns an `O`" — any function or
lambda with that shape already qualifies, with zero implementing required.

### Who calls whom

```
main()
 ├── builds stage1 = NewPipeline(removeAlphabets)
 ├── builds stage2 = AddHandler(stage1, removeDigits)
 ├── builds stage3 = AddHandler(stage2, toCharArray)
 └── calls stage3.Execute(input)
```

### What actually happens, in order, when you run it

1. `main()` wraps `removeAlphabets` in `NewPipeline`, producing `stage1` — a pipeline that can only remove letters.
2. `main()` calls `AddHandler(stage1, removeDigits)`, producing `stage2` — a *new* pipeline whose one function is "run `stage1`, then run `removeDigits` on the result." Nothing has run yet.
3. `main()` calls `AddHandler(stage2, toCharArray)`, producing `stage3` — nested one layer deeper again.
4. `main()` calls `stage3.Execute("GoYankees123!")`. This is the single moment that fires all 3 stages, nested-call style, same as the Java version: `"123!"` → `"!"` → `['!']`.
5. The result is printed.

### Method-by-method: who calls what

| File | Function | Called by | It calls |
|---|---|---|---|
| `main.go` | `removeAlphabets(input)` | wrapped into a `Handler` inside `main()` | nothing external — pure loop over the string |
| `main.go` | `removeDigits(input)` | same | nothing external |
| `main.go` | `toCharArray(input)` | same | nothing external |
| `pipeline.go` | `NewPipeline(h)` | `main()` | nothing — just stores `h` |
| `pipeline.go` | `AddHandler(p, next)` | `main()`, once per extra stage | builds a new function that calls `p.run`, then `next` |
| `pipeline.go` | `(Pipeline).Execute(input)` | `main()` | the pipeline's one (possibly nested) function |
| `main.go` | `main()` | program entry point | `NewPipeline`, `AddHandler` (twice), `Execute` |

**One real Go-vs-Java difference worth knowing for an interview**: in Java,
`addHandler` is an *instance method* on `Pipeline` that introduces its own
extra generic type (`<K> Pipeline<I,K> addHandler(Handler<O,K> next)`),
which is why Java can chain it fluently as `.addHandler(...).addHandler(...)`.
Go **does not allow a method to introduce a type parameter beyond the ones
its receiver already has** — so `AddHandler` in the Go version has to be a
standalone function, not a method, and you assign each result to a new
variable (`stage2 := AddHandler(stage1, ...)`) instead of chaining dots.
Same idea, different shape, because of a real language-level constraint.

---

## Side-by-side: same idea, different word

| What it does | Java | Go |
|---|---|---|
| The "one-method contract" every stage follows | `Handler<I,O>` interface | `Handler[I,O]` — just a named function type, no interface needed |
| A single stage | a lambda or method reference implementing `Handler` | a plain function, or a value of type `Handler[I,O]` |
| Chaining one more stage on | `pipeline.addHandler(next)` — an instance method | `AddHandler(pipeline, next)` — a standalone function |
| Running the whole chain | `pipeline.execute(input)` | `pipeline.Execute(input)` |
| What actually nests the stages together | a lambda capturing the previous handler | a closure capturing the previous pipeline's `run` function |

---

## Pseudocode you can write from scratch in an interview

```
define Stage as: a function that takes one value, returns one value

define Pipeline holding: one Stage (call it "run")

function newPipeline(stage):
    return Pipeline{ run: stage }

function addStage(pipeline, nextStage):
    oldRun = pipeline.run
    return Pipeline{ run: (input) -> nextStage(oldRun(input)) }
    # this is the whole trick: wrap the old stage and the new stage
    # together into ONE new function, don't run anything yet

function execute(pipeline, input):
    return pipeline.run(input)   # this single call cascades through every
                                  # stage that was ever added, in order

# usage
p = newPipeline(stage1)
p = addStage(p, stage2)
p = addStage(p, stage3)
result = execute(p, firstInput)
```

The one line worth memorizing: **"adding a stage doesn't run anything — it
builds a bigger function that will run everything, later, all at once,
when you finally call execute."**

---

## Likely interview questions, answered short

**What problem does this solve?**
Turning one big transformation into several small, independently
understandable and testable steps, chained together in a fixed order.

**Is Pipeline a concurrency pattern?**
Not inherently — the version here runs every stage on the same thread/
goroutine, one after another. It *can* be made concurrent (each stage on
its own thread/goroutine, connected by queues — like a Unix
`cmd1 | cmd2 | cmd3`), but that's a variation, not the core idea. The core
idea is ordered, staged, type-checked transformation.

**How is a stage prevented from being plugged in at the wrong point in the
chain?**
By the type parameters. `addHandler`/`AddHandler` requires the new stage's
input type to match the *current* pipeline's output type. If stage 2
expects an `int` but stage 1 outputs a `String`, it's a compile error, not
a runtime surprise.

**Why did the Java version originally have a separate class per stage, and
why was that overkill?**
Because `Handler` only has one method, any class implementing it is doing
nothing but supplying that one method — which is exactly what a lambda
already is. A whole file/class per stage added navigation overhead without
adding any capability a lambda didn't already have.

**How is this different from Chain of Responsibility?**
Chain of Responsibility passes the *same* type along a chain, and any link
can decide to stop the chain early (e.g., "I handled it, don't pass it
on"). Pipeline always runs every stage, in order, and stages are allowed to
change the type at each step (`String` → `String` → `char[]` here) — there's
no "stop early" concept.
