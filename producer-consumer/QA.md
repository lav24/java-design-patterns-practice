# Producer-Consumer — Interview Q&A

Real questions, how to approach each one out loud, and a model answer.
Pair this with `EXECUTION_FLOW.md` for the mechanics.

---

**Q: What is the Producer-Consumer pattern, and what problem does it
solve?**

*How to answer:* Definition, then the specific benefit — decoupling rate,
not just decoupling code.

> One or more producers keep generating work, one or more consumers keep
> processing it, connected only through a shared bounded queue — neither
> side knows the other exists. The real problem it solves is letting two
> things that run at different, possibly varying speeds cooperate without
> one blocking the other's code path directly, and without unbounded memory
> growth if one side temporarily outpaces the other.

---

**Q: What is backpressure, and how does this pattern provide it?**

*How to answer:* Define the term plainly first — some interviewers use
"backpressure" expecting you to know the word; don't dodge it, name it and
explain it in one sentence.

> Backpressure is a fast producer being automatically slowed down to match
> a slower consumer, instead of piling up unbounded work in memory. Here it
> comes directly from the queue's fixed capacity: `put()` (or a channel
> send in Go) blocks once the queue is full, so a producer physically
> cannot get more than N items ahead of the consumers — the queue itself
> enforces the slowdown, no extra code needed.

---

**Q: Why use a `BlockingQueue` instead of a plain `Queue` with manual
`wait()`/`notify()`?**

*How to answer:* Name the two things `BlockingQueue` gives you for free
that hand-rolled wait/notify makes you get right yourself.

> `BlockingQueue` already implements exactly the "pause producer when full,
> pause consumer when empty, wake the right side up when state changes"
> logic — correctly, including edge cases around spurious wakeups that are
> easy to get subtly wrong by hand. Writing it yourself with `wait()`/
> `notify()` means re-deriving battle-tested concurrency code that the
> standard library already ships, for no benefit.

---

**Q: What happens if the queue is unbounded instead of capped at a fixed
size?**

*How to answer:* Name the concrete failure mode, not just "it's bad
practice."

> If producers are ever faster than consumers for a sustained period, an
> unbounded queue lets unprocessed items accumulate without limit — that's
> unbounded memory growth, and eventually an `OutOfMemoryError` or the
> equivalent resource exhaustion. Capping the queue converts "the process
> eventually crashes" into "producers pause," which is almost always the
> outcome you actually want.

---

**Q: How do you cleanly stop all the producers and consumers once you're
done?**

*How to answer:* Describe what the demo code actually does, and be honest
that it's a coarser mechanism in Java than in Go — that contrast is worth
volunteering.

> In the Java version here, each worker runs `while (true)`, so stopping
> means calling `executorService.shutdownNow()`, which interrupts every
> running thread individually — the interrupt breaks whichever blocking
> call (`put`/`take`) that thread happened to be parked in. In Go, a
> `done` channel gets closed once, and every goroutine's `select` picks
> that up on its own — one call, every goroutine notices, instead of
> interrupting each one separately. Go's version is a cleaner one-to-many
> broadcast for this specific case.

---

**Q: What issues could come up with multiple producers and multiple
consumers sharing one queue, beyond just correctness?**

*How to answer:* This is asking for something a bit deeper than the basic
mechanics — fairness and contention are the right things to raise.

> Fairness: `BlockingQueue` doesn't guarantee items are consumed in strict
> arrival order across multiple producers, and it doesn't guarantee which
> waiting consumer gets woken up first when an item arrives — usually fine,
> but worth knowing if ordering matters for your use case. There's also
> contention: with many producers and consumers all hitting the same
> queue, the queue itself can become a bottleneck; at high scale you'd look
> at sharding across multiple queues or a purpose-built message broker
> instead.

---

**Q: How is this different from Master-Worker?**

*How to answer:* Same crisp one-sentence contrast as the Master-Worker
QA — consistency across your answers here signals you actually understand
both, not just memorized two separate blurbs.

> Producer-Consumer has no concept of "done" — it's meant to run
> indefinitely, taking in new work for as long as the program is alive.
> Master-Worker has a fixed, known batch — someone waits for every piece to
> finish, then the whole thing is over.

---

**Q: Where have you seen this pattern used in real systems?**

*How to answer:* Give concrete, named examples — this is an easy question
to answer well if you connect it to infrastructure you actually know about.

> Message queues and brokers like Kafka or RabbitMQ are Producer-Consumer
> at the infrastructure level — services produce messages, other services
> consume them, decoupled by the broker. Thread pools are also a variant:
> a work queue holds submitted tasks (produced by application code), and
> pool threads consume and execute them. Logging frameworks often use it
> internally too — application threads produce log lines, a dedicated
> background thread consumes and writes them, so logging I/O never blocks
> the calling thread.

---

**Q: Why does Go need one fewer file than Java for this pattern?**

*How to answer:* Point to the specific language feature responsible.

> Java has no built-in "queue that automatically pauses producers and
> consumers," so the project wraps one (`BlockingQueue`) in a small
> `ItemQueue` class. Go has that exact behavior — bounded capacity, blocking
> send/receive — built directly into the language as a channel, so there's
> no wrapper class needed at all; `make(chan Item, 5)` *is* the queue.
