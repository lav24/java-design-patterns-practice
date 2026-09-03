# Pipeline — Interview Q&A

Real questions, how to approach each one out loud, and a model answer.
Pair this with `EXECUTION_FLOW.md` for the mechanics.

---

**Q: What is the Pipeline pattern?**

*How to answer:* Definition, then immediately clarify the common
misconception — that this is a concurrency pattern. Volunteering that
correction unprompted signals you actually understand it rather than
pattern-matching the name.

> A value passes through an ordered sequence of independent transformation
> stages, each stage's output becoming the next stage's input — like an
> assembly line. It's easy to assume this is a concurrency pattern because
> it sits next to Master-Worker and Producer-Consumer in most pattern
> catalogs, but the version I built runs every stage sequentially, on one
> thread. The core idea — staged, ordered, type-checked composition — is
> actually independent of whether the stages run concurrently or not.

---

**Q: Can Pipeline be made concurrent? What would that look like?**

*How to answer:* Yes, describe the shape briefly — this shows you
understand the pattern is a spectrum, not one fixed implementation.

> Yes — each stage could run on its own thread or goroutine, connected to
> the next by a queue or channel, similar to a Unix `cmd1 | cmd2 | cmd3`
> pipe. That gets you throughput benefits when stages have very different
> costs, since a slow stage doesn't block the whole chain from accepting new
> input at stage one. It's a valid variation, but a different implementation
> choice from the type-safe function-composition version I built.

---

**Q: How does `Pipeline.addHandler` achieve compile-time type safety across
stages?**

*How to answer:* Walk through the generic signature, and be precise about
*what* gets checked and *when* — this is the mechanically interesting part
of the pattern.

> `addHandler` is generic over a new type parameter — `<K> Pipeline<I, K>
> addHandler(Handler<O, K> newHandler)` — where `O` is the *current*
> pipeline's output type. That means the compiler requires the new stage's
> input type to exactly match the previous stage's output type at the call
> site, or it won't compile. If I tried to plug a stage expecting an `int`
> after one that outputs a `String`, that's a compile error, not something
> that surfaces as a runtime `ClassCastException` later.

---

**Q: When you call `addHandler` three times, does anything actually run at
that point?**

*How to answer:* No — say so directly, then explain what actually gets
built instead. This is testing whether you understand laziness versus
eager evaluation.

> No — `addHandler` only builds a new, more deeply nested function; it
> returns immediately without invoking anything. Each call wraps the
> previous handler and the new one together into a single lambda that says
> "run the old chain, then run the new stage on its result." Nothing
> actually executes until `execute(input)` is called, at which point that
> one (possibly deeply nested) function runs all the way through in a
> single call.

---

**Q: This module used to have a separate named class for every stage. Why
did you replace them with lambdas, and how did you justify that it was
safe?**

*How to answer:* Name the language feature that makes it possible
(single-method interface), then be specific about what you checked before
making the change — showing verification discipline matters here.

> `Handler<I, O>` only declares one method, `process`, which makes it a
> functional interface — any class implementing it is really just
> supplying that one method, exactly what a lambda already is. So each of
> those three classes was a whole file adding navigation overhead without
> adding any capability. Before deleting them, I checked exactly what the
> existing tests referenced — one test constructed the handler classes
> directly by name, so I moved the lambdas into small static factory
> methods on `App` and updated that one test to call those factories
> instead, keeping the exact same input/output behavior.

---

**Q: How is Pipeline different from Chain of Responsibility?**

*How to answer:* Name the two concrete differences — type flexibility and
early-exit — rather than a vague "they're similar but different."

> Chain of Responsibility passes the *same* type along the whole chain, and
> any link in the chain can choose to stop it early — "I handled this,
> don't pass it further." Pipeline always runs every stage in order with no
> early exit, and stages are explicitly allowed to change the type at each
> step — this example goes `String` → `String` → `char[]`. Different
> intents: Chain of Responsibility is about "who handles this," Pipeline is
> about "transform this through a fixed sequence of steps."

---

**Q: In the Go version, why can't `addHandler` be a method on `Pipeline`
the way it is in Java?**

*How to answer:* This is a real, specific language constraint — naming it
precisely is worth more than a vague "Go generics are different."

> Go doesn't allow a method to introduce a type parameter beyond the ones
> its receiver already has. Java's `addHandler` needs its own type
> parameter `<K>` for the new stage's output type, in addition to the
> `Pipeline<I, O>`'s existing `I` and `O` — that's legal for a Java instance
> method but not for a Go method. So in Go, `AddHandler` has to be a
> standalone generic function taking the pipeline as its first argument
> instead of a method you can chain with dots — same idea, different shape,
> because of that one language rule.

---

**Q: Does Go need a `Handler` interface the way Java does?**

*How to answer:* No — and explain precisely why, since this is a genuinely
interesting language contrast to raise unprompted.

> No — in Go, `Handler[I, O]` is just a named function type, not an
> interface with a method to implement. Go functions are already values you
> can pass around directly, so any function or lambda with the matching
> shape (`func(I) O`) already satisfies it with nothing to implement at
> all. Java needed an actual interface because a lambda in Java has to be
> matched against some functional interface's single abstract method before
> the compiler will accept it — Go has no such requirement.
