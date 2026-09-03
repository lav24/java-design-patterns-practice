# Factory Method — Interview Q&A

Real questions, how to approach each one out loud, and a model answer.
Pair this with `EXECUTION_FLOW.md` for the mechanics.

---

**Q: What is the Factory Method pattern?**

*How to answer:* Definition first, then the specific mechanism
(polymorphism), since that's what separates it from "any method that
returns an object."

> It lets a class defer the decision of which concrete class to instantiate
> to its subclasses (or, in an interface-based version, to whichever
> concrete implementation the caller is holding). The caller only ever
> interacts through a product interface and a creator interface — it never
> writes the concrete class name itself.

---

**Q: Is any method that returns a new object via `new` a "factory
method"?**

*How to answer:* No — this is a common misconception worth correcting
directly, since it's an easy trap to fall into if you've only seen the name
without the mechanism.

> Not in the GoF sense. A static utility method like
> `WeaponFactory.create(String type) { if (type.equals("elf")) return new
> ElfWeapon(); ... }` is sometimes called a "Simple Factory," but it has no
> subclassing or polymorphic dispatch — it's just conditional logic in one
> method. The actual Factory Method pattern specifically relies on
> different concrete creators (subclasses or separate implementations)
> overriding a shared method to each produce their own product type — the
> decision is made by *which object you're holding*, not by an if/else
> branch inside one method.

---

**Q: In this codebase's example, `ElfBlacksmith` only ever produces
`ElfWeapon`s and `OrcBlacksmith` only ever produces `OrcWeapon`s regardless
of the `WeaponType` argument. So what is the `WeaponType` parameter
actually for?**

*How to answer:* This is worth answering precisely, because it's easy to
conflate the two things happening. Name both axes explicitly.

> There are two independent decisions in this example, and only one of them
> is Factory Method. Which *race* of weapon comes back — elf or orc — is
> decided by which concrete `Blacksmith` you're holding, and that's the
> actual polymorphic dispatch the pattern is about. Which *weapon type*
> comes back — spear, axe, short sword — is decided by the `WeaponType`
> parameter, but that's just a lookup into a pre-built map; a given
> blacksmith never changes which race it produces based on that parameter.
> If you removed the parameter entirely and gave each blacksmith a single
> hardcoded weapon, you'd still have a complete, arguably clearer
> demonstration of Factory Method — the parameter demonstrates a separate,
> smaller idea (a factory method producing one of several related products
> within a family) bolted onto the same example.

---

**Q: How is Factory Method different from Abstract Factory?**

*How to answer:* Name the scope difference precisely — one method vs. a
family of related methods/objects.

> Factory Method is about **one** product, created via **one** overridable
> method. Abstract Factory is about producing a **whole family** of related
> products that need to be used together — one factory interface with
> *multiple* creation methods, and swapping the whole factory implementation
> swaps every product in the family consistently. Abstract Factory is
> often implemented *using* several Factory Methods internally, one per
> product in the family.

---

**Q: Why does the caller in this example declare its variables using the
interface types (`Blacksmith`, `Weapon`) instead of the concrete classes?**

*How to answer:* This is really asking "do you understand what the pattern
actually buys you," not just "can you name the types."

> Because the entire benefit of the pattern is that the caller's code
> doesn't change when you swap which concrete factory you're using. If
> `App` declared `OrcBlacksmith blacksmith = new OrcBlacksmith();` instead
> of `Blacksmith blacksmith = ...`, then switching to `ElfBlacksmith` later
> would require changing that declaration too — you'd have lost the whole
> point. Declaring against the interface means one line (the construction)
> is the only place the concrete type is ever named.

---

**Q: What are the trade-offs of using this pattern versus just calling
constructors directly?**

*How to answer:* Give one benefit and one honest cost — a one-sided answer
reads as memorized rather than understood.

> The benefit is decoupling — calling code depends only on an interface, so
> new product types can be added by writing a new concrete creator, without
> touching any existing calling code (this is the Open/Closed Principle in
> practice). The cost is indirection — for a simple case with only one or
> two product types that will never grow, introducing a creator interface
> and multiple concrete factory classes can be more ceremony than the
> problem needs; I'd reach for this pattern when I actually expect the set
> of product types to grow or vary by configuration, not preemptively.

---

**Q: How does Go's implicit interface satisfaction change how you'd
implement this pattern compared to Java?**

*How to answer:* Name the specific mechanical difference — no `implements`
keyword — and what that implies practically.

> In Java, a class has to explicitly declare `implements Blacksmith`. In
> Go, any type automatically satisfies an interface just by having methods
> with matching signatures — there's no keyword linking them at all. Practically,
> that means in Go you could pass an existing, unrelated type into code
> expecting a `Blacksmith` as long as it happens to have a matching
> `ManufactureWeapon` method, with zero changes to that type — you can't
> retrofit an interface onto a class after the fact in Java without editing
> the class itself.

---

**Q: In the Go port, why are `ElfBlacksmith`/`OrcBlacksmith` defined with
value receivers instead of pointer receivers?**

*How to answer:* State the rule, then confirm it against this specific
case — this shows you're applying a principle, not reciting a fact.

> The rule of thumb is: use a pointer receiver when a method needs to
> modify the actual object, or when the struct is large enough that copying
> it on every call would be wasteful. Neither applies here — nothing about
> a `Blacksmith` ever changes after it's constructed, and the struct is
> just one `map` field, which is already cheap to copy since a map value is
> a small header pointing at the real data. So value receivers are simpler
> and correct; reaching for a pointer here would just be defaulting to
> reference semantics out of habit rather than need.
