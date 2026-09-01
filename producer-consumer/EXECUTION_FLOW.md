# Producer-Consumer — Who Calls What, and Why

This is a study reference, not the project's official docs page (that's `README.md`).
Goal: know exactly which file talks to which file, in what order, in both Java and Go,
well enough to draw it on a whiteboard and write it from memory.

## The big picture, in one paragraph

There's one shared "waiting line" (a queue) with a limited number of spots.
Some workers ("producers") keep adding new items to the line. Other workers
("consumers") keep taking items off the line. If the line is full, a producer
has to wait for a spot to open up. If the line is empty, a consumer has to
wait for something to arrive. Producers and consumers never talk to each
other directly — they only ever talk to the line.

---

## Java version

### The files and what each one is

| File | What it is |
|---|---|
| `Item.java` | The thing being passed around. Just data: which producer made it, and an id number. |
| `ItemQueue.java` | The waiting line itself. Wraps Java's built-in `BlockingQueue`, capped at 5 spots. |
| `Producer.java` | A worker that keeps making new `Item`s and placing them on the line. |
| `Consumer.java` | A worker that keeps taking `Item`s off the line. |
| `App.java` | The program's starting point. Builds everything and starts the workers running. |

### Who holds a reference to whom (the relationships)

```
App
 ├── creates ONE  ItemQueue                         (the shared line)
 ├── creates TWO  Producer objects  ─── each one is handed the SAME ItemQueue
 └── creates THREE Consumer objects ─── each one is handed the SAME ItemQueue

Producer  ---holds a reference to--->  ItemQueue
Consumer  ---holds a reference to--->  ItemQueue

(Producers and Consumers do NOT hold a reference to each other. The
 ItemQueue is the only thing connecting them.)
```

### What actually happens, in order, when you run it

1. `App.main()` starts. It creates one `ItemQueue`.
2. It creates a pool of 5 worker slots (`Executors.newFixedThreadPool(5)`) — think of this as 5 empty desks that will each run one worker.
3. It creates 2 `Producer` objects (`Producer_0`, `Producer_1`), each one built with the same `ItemQueue`. Each is handed to one desk, and told: "run `produce()` over and over, forever."
4. It creates 3 `Consumer` objects (`Consumer_0`, `Consumer_1`, `Consumer_2`), each built with the same `ItemQueue`. Each is handed to one desk, and told: "run `consume()` over and over, forever."
5. From this point, all 5 workers run independently and simultaneously — nobody is waiting for anybody else to go first.
6. Whenever a `Producer`'s `produce()` runs: it builds a new `Item`, then calls `ItemQueue.put(item)`. If the line already has 5 items sitting in it, this call pauses that one producer until a consumer removes something. Then it sleeps for a random moment (simulating "it takes time to make the next item") and loops back to step 6.
7. Whenever a `Consumer`'s `consume()` runs: it calls `ItemQueue.take()`. If the line is empty, this call pauses that one consumer until a producer adds something. Once it gets an `Item`, it prints who made it and loops back to step 7.
8. This continues until `App.main()` decides to stop everyone (after 10 seconds in the demo), at which point it forces all 5 workers to stop mid-wait.

### Method-by-method: who calls what

| File | Method | Called by | It calls |
|---|---|---|---|
| `Item.java` | constructor `Item(producer, id)` | `Producer.produce()` | nothing — it's just data |
| `ItemQueue.java` | `put(item)` | `Producer.produce()` | the built-in queue's `put()` (Java's standard library does the actual waiting) |
| `ItemQueue.java` | `take()` | `Consumer.consume()` | the built-in queue's `take()` |
| `Producer.java` | `produce()` | `App.main()`, inside a loop that never stops | `new Item(...)`, `queue.put(item)`, `Thread.sleep(...)` |
| `Consumer.java` | `consume()` | `App.main()`, inside a loop that never stops | `queue.take()`, then a print statement |
| `App.java` | `main()` | the program starts here — nobody calls it, the computer does | creates `ItemQueue`, `Producer`, `Consumer`; starts the worker pool; eventually stops it |

---

## Go version

### The files and what each one is

| File | What it is |
|---|---|
| `producerconsumer.go` | Defines `Item` (the data) and the two worker functions: `produce` and `consume`. |
| `main.go` | The program's starting point. Builds the shared line and launches the workers. |

Go doesn't need a separate "queue class" file — a **channel** (a built-in Go
feature) already behaves exactly like the Java `ItemQueue`: it has a fixed
capacity, and reading/writing it automatically pauses a worker when it's
empty/full. So there's one less file than Java, because the language itself
provides what `ItemQueue.java` had to build by hand.

### Who holds a reference to whom

```
main()
 ├── creates ONE   channel (the shared line, holds up to 5 items)
 ├── creates ONE   "done" channel (used only to say "stop now")
 ├── launches TWO  produce() workers ─── each given the SAME channel + done signal
 └── launches THREE consume() workers ─── each given the SAME channel + done signal
```

### What actually happens, in order, when you run it

1. `main()` starts. It creates the shared channel (capacity 5) and a second, empty "stop signal" channel.
2. It launches 2 `produce` workers (`go produce(...)`) — each one starts running immediately and independently, sharing the same channel.
3. It launches 3 `consume` workers (`go consume(...)`) the same way.
4. Every `produce` worker loops: build an `Item`, then try to place it on the channel. If the channel already holds 5 items, this pauses that one worker until a consumer removes something. Then it sleeps briefly and loops again.
5. Every `consume` worker loops: try to take an item off the channel. If it's empty, this pauses that one worker until a producer adds something. Once it gets an item, it prints it and loops again.
6. After some time, `main()` closes the "stop signal" channel. The instant it does, every worker's next check of that signal succeeds, and all of them exit — one action stops every worker at once (Java needs to interrupt each worker individually to get the same effect).

### Method-by-method: who calls what

| File | Function | Called by | It calls |
|---|---|---|---|
| `producerconsumer.go` | `Item{...}` (creating one) | `produce()` | nothing — it's just data |
| `producerconsumer.go` | `produce(name, queue, done)` | started directly by `main()` via `go produce(...)` | sends on `queue`, sleeps, checks `done` |
| `producerconsumer.go` | `consume(name, queue, done)` | started directly by `main()` via `go consume(...)` | receives from `queue`, prints, checks `done` |
| `main.go` | `main()` | the program starts here | creates the channels, starts `produce`/`consume` workers, later closes `done` |

---

## Side-by-side: same idea, different word

| What it does | Java | Go |
|---|---|---|
| The shared line itself | `ItemQueue` class wrapping `BlockingQueue` | a channel: `make(chan Item, 5)` |
| Limiting how many items can wait | `new LinkedBlockingQueue<>(5)` | the `5` in `make(chan Item, 5)` |
| An independent worker | a thread from the thread pool running a loop | a goroutine (`go produce(...)`) |
| Pausing a worker when the line is full/empty | happens automatically inside `put()`/`take()` | happens automatically on channel send/receive |
| Telling every worker to stop at once | not automatic — each thread is interrupted individually | `close(done)` — one call, every worker notices |
| The item being passed around | `Item` (a record) | `Item` (a struct) |

---

## Pseudocode you can write from scratch in an interview

This is language-neutral — the shape is identical whether you write it in
Java, Go, Python, or anything else.

```
define Item: { producerName, id }

define SharedLine with a maximum size N:
    function put(item):
        if line is full: wait here until there's room
        add item to the line
        if anyone is waiting to take an item, wake one of them up

    function take():
        if line is empty: wait here until something arrives
        remove and return the oldest item
        if anyone is waiting to put an item, wake one of them up

function producerJob(name, line):
    counter = 0
    repeat forever:
        item = Item(name, counter)
        line.put(item)            # this line may pause here
        counter = counter + 1
        wait a random short amount of time

function consumerJob(name, line):
    repeat forever:
        item = line.take()        # this line may pause here
        print "consumed", item, "by", name

function main():
    line = SharedLine(maxSize = 5)
    start 2 independent workers, each running producerJob(name, line)
    start 3 independent workers, each running consumerJob(name, line)
    # they all run at the same time, forever, until something stops them
```

If asked to implement the shared line's `put`/`take` yourself (without a
built-in blocking queue/channel), the plain-English version of the waiting
logic is: *"keep a count of how many items are currently in the line; a
producer checks that count before adding and pauses if it's already at the
max; a consumer checks that count before removing and pauses if it's zero;
whenever either one changes the count, it notifies whoever might be
waiting."* That's the textbook wait/notify implementation interviewers
sometimes want you to write by hand instead of using a library.

---

## Likely interview questions, answered short

**What problem does this solve?**
Letting two things that run at different speeds work together without one
having to wait around doing nothing for the other, and without them needing
to know anything about each other.

**What happens if producers are faster than consumers?**
The line fills up (hits its max of 5), and producers start pausing on
`put`/send until consumers catch up. This automatic slowing-down is called
backpressure.

**What happens if consumers are faster than producers?**
Consumers pause on `take`/receive, waiting for the next item. No wasted
CPU spinning in a loop checking "is there anything yet?" — the pause is
free.

**Why cap the line at a fixed size instead of letting it grow forever?**
An uncapped line lets a fast producer pile up unlimited items in memory if
consumers ever fall behind — that can crash the program. A cap forces the
slowdown to happen instead.

**How is this different from Master-Worker?**
Master-Worker has a fixed, known amount of work — someone waits until every
piece is done, then the whole thing finishes. Producer-Consumer has no
concept of "done" — it's meant to run indefinitely, processing whatever
shows up, for as long as the program is alive.

**Why does Go need one less file than Java?**
Java doesn't have a built-in "queue that automatically pauses producers and
consumers" as a language feature, so the project writes `ItemQueue.java` as
a thin wrapper around one (`BlockingQueue`) from Java's standard library. Go
has that exact behavior built directly into the language as a channel, so
there's nothing left to wrap.
