# Singleton — Interview Q&A

Real questions, how to approach each one out loud, and a model answer.
Pair this with `EXECUTION_FLOW.md` for the mechanics.

---

**Q: What is the Singleton pattern, and when would you actually use it?**

*How to answer:* Give the one-line definition, then immediately name a
concrete real use case — interviewers want to know you understand *why*,
not just *what*.

> A Singleton guarantees a class has exactly one instance for the life of
> the program, with one shared access point to it. I'd reach for it for
> things that are genuinely singular by nature and expensive to duplicate —
> a logging framework, a configuration object loaded once from disk, a
> connection pool. I'd be cautious using it just to avoid passing a
> parameter around, though — that's usually a sign dependency injection is
> the better tool.

---

**Q: Walk me through the trade-offs between the different ways to implement
it.**

*How to answer:* Go in order of complexity, and for each one name the ONE
thing it trades away versus the previous one — that structure makes it
sound rehearsed-but-understood rather than memorized.

> Eager initialization is the simplest — thread-safe for free, because the
> JVM only initializes a class once — but it builds the object even if the
> program never ends up using it. A plain `synchronized` getter fixes the
> "build only when needed" problem but pays a locking cost on literally
> every call, forever, even after the instance already exists.
> Double-checked locking fixes that — cheap unlocked check first, only
> locks the first few times — but it's easy to get wrong; it needs
> `volatile` or it's actually broken under certain memory-reordering
> conditions. The holder idiom gets you lazy, thread-safe, AND zero runtime
> locking, by leaning on the JVM's class-loading guarantees instead of
> hand-written synchronization — that's usually my default. Enum is the
> outlier: it's the only one immune to reflection and deserialization
> attacks, so it's the answer when "bulletproof" actually matters.

---

**Q: Why does double-checked locking need the field to be `volatile`? What
breaks if it isn't?**

*How to answer:* This is a filter question — a lot of candidates can name
"double-checked locking" but can't explain why `volatile` specifically
matters. Explain what "reordering" means concretely, not just as a buzzword.

> Building an object isn't one atomic step — the JVM can, in principle,
> make the reference visible to another thread before the constructor has
> fully finished running on it, because instruction reordering is legal in
> the absence of the `volatile` field's memory barrier. So without
> `volatile`, a second thread's outer null-check could see a non-null
> reference and start using an object whose fields haven't been set yet — a
> half-built object. `volatile` prevents that specific reordering, so
> "reference is visible" and "construction is finished" always happen in
> the order you'd expect. This was actually a well-known broken pattern
> pre-Java-5, before the memory model was fixed to make `volatile` strong
> enough to guarantee this.

---

**Q: How would you break someone's Singleton using reflection?**

*How to answer:* Show, don't just assert — walk through the actual API
call, then be ready for the natural follow-up ("so how do you defend
against it?").

> `Constructor.setAccessible(true)` bypasses Java's access-control check
> entirely, letting you call a `private` constructor from outside the
> class: `var ctor = MyClass.class.getDeclaredConstructor(); ctor.setAccessible(true); var forged = ctor.newInstance();`
> `private` is a compile-time rule the compiler enforces — reflection
> operates below that layer and doesn't ask the compiler at all.

**Follow-up you should expect: "So how do you actually defend against
it?"**

> Every implementation in this codebase adds a guard in the constructor:
> `if (instance != null) throw new IllegalStateException(...)`. I actually
> tested this — the interesting part is *why* it works: calling a
> constructor forces the JVM to finish that class's static initialization
> first, whether the call is reflective or not. So the legitimate instance
> always gets built — via a normal internal call, where the guard correctly
> sees `null` — before the attacker's reflective call ever runs; by the time
> that runs, the guard sees non-null and throws. I actually verified this
> against an *unguarded* version too — with `private` and nothing else, the
> attack fully succeeds and creates a second, independent object. So it's
> the guard code doing the work, not the `private` keyword by itself. Enum
> is the one exception — the JVM refuses to reflectively construct an enum
> at all, unconditionally, with a hardcoded check, no guard code required.

---

**Q: How can a Singleton be broken by deserialization, and how do you fix
it?**

*How to answer:* Name the root cause first (deserialization skips
constructors), then the fix, then the one implementation that needs
neither.

> Deserialization doesn't call the constructor at all — it allocates a raw
> object and populates it directly from the byte stream, so any
> "already initialized" guard in the constructor is irrelevant here; it
> never runs. If a Singleton implements `Serializable` with no other
> precaution, deserializing a previously-saved copy hands you back a
> second, distinct object — I verified this concretely: identical class,
> `original == rebuilt` came back `false`. The fix is a `readResolve()`
> method — deserialization calls it automatically right after building the
> raw object and substitutes whatever it returns for the caller. Return the
> real singleton instance from it, and the caller never sees the throwaway
> object that got built internally. Enum needs none of this — Java's
> serialization spec serializes an enum constant by name and resolves it
> via `Enum.valueOf` on read, never reconstructing it from scratch.

---

**Q: Is Singleton considered an anti-pattern? Why or why not?**

*How to answer:* Don't just say "yes" or "no" — show you can see both
sides, then state your actual position.

> It's controversial for good reasons: it's global mutable state, it hides
> a dependency instead of making it an explicit parameter/constructor
> argument, and it makes unit testing harder because you can't easily swap
> in a mock or a fresh instance between tests. That said, the *idea* — "one
> shared instance, one access point" — isn't wrong; it's the hand-rolled,
> globally-reachable implementation that's the problem. In most modern
> codebases I'd let a framework (like Spring's default singleton-scoped
> bean) manage that lifecycle via dependency injection instead of writing
> `getInstance()` static methods myself — same guarantee, but injectable and
> mockable.

---

**Q: How would you implement this in a language without classes, like Go?**

*How to answer:* Lead with what maps cleanly, then be explicit about what
genuinely doesn't translate — that's the part that shows real understanding
instead of pattern-matching syntax across languages.

> Go doesn't have constructors, so there's no direct equivalent of "private
> constructor plus a guard." The eager case maps cleanly: a package-level
> variable initializes exactly once, before `main()` runs, for free — same
> guarantee as Java's eager static field. For the lazy case, Go's standard
> library has `sync.Once`, which basically *is* double-checked locking,
> already written and race-tested — `once.Do(func(){...})` runs its
> function exactly once no matter how many goroutines call it concurrently,
> and every other caller blocks until that run finishes. What doesn't
> translate: there's no reflection-attack surface the way Java has, because
> there's no constructor to bypass in the first place — the whole guarantee
> in Go rests on the type being unexported (lowercase), which is a package
> boundary, not a runtime-enforced rule. And Go has nothing like Java's enum
> singleton — no language construct the runtime specially protects against
> forgery.

---

**Q: `InitializingOnDemandHolderIdiom` and `BillPughImplementation` look
almost identical — what's actually different between them?**

*How to answer:* This is a trap question in this specific codebase — the
correct answer is "nothing, structurally."

> Nothing meaningfully different — they're the same technique, described in
> two different sources under two different names (Bill Pugh's write-up,
> and the more generic "initialization-on-demand holder idiom" name). Both
> rely on the exact same mechanism: a static nested class isn't loaded until
> it's first referenced, so wrapping the instance field in one gives you
> lazy loading, and the JVM's own once-only class-initialization guarantee
> gives you thread safety, with no explicit `synchronized` or `volatile`
> anywhere.
