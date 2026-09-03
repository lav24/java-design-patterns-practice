# Factory Method — Who Calls What, and Why

This is a study reference, not the project's official docs page (that's `README.md`).
Goal: know exactly which file talks to which file, in what order, in both Java and Go,
well enough to draw it on a whiteboard and write it from memory.

## The big picture, in one paragraph

Code that needs an object shouldn't have to know or write the exact
concrete class it wants — it should ask through an interface, and let
something else (a specific factory) decide which concrete class actually
gets built. Swap which factory you're holding, and every line of calling
code produces a different concrete type without a single other line
changing.

## The one thing worth separating out before anything else

This example has **two independent decisions** happening, and only one of
them is what Factory Method is actually about:
- **Which race (elf/orc)** — decided by which `Blacksmith` you construct. **This is Factory Method** — polymorphic dispatch across `Blacksmith` subclasses/implementations.
- **Which weapon type (spear/axe/sword)** — decided by a parameter passed into `manufactureWeapon`. **This is not Factory Method** — it's just a lookup, no polymorphism involved; a given `Blacksmith` always returns the same race of weapon no matter what type you ask for.

---

## Java version

### The files and what each one is

| File | What it is |
|---|---|
| `Weapon.java` | The product interface. |
| `WeaponType.java` | An enum of weapon kinds — data, not part of the factory mechanism. |
| `ElfWeapon.java` / `OrcWeapon.java` | The two concrete products. |
| `Blacksmith.java` | The creator interface — declares the factory method. |
| `ElfBlacksmith.java` / `OrcBlacksmith.java` | The two concrete factories, each hardcoded to one product race. |
| `App.java` | Demonstrates swapping which factory is held. |

### Who calls whom (the relationships)

```
App
 └── holds a Blacksmith reference (compile-time type is the INTERFACE)
        ├── could point to: ElfBlacksmith  ──produces──>  ElfWeapon
        └── could point to: OrcBlacksmith  ──produces──>  OrcWeapon

App never mentions ElfBlacksmith/OrcBlacksmith/ElfWeapon/OrcWeapon by name
except at the ONE line where it decides which to construct.
```

### What actually happens, in order, when you run it

1. `App.main()` constructs `new OrcBlacksmith()` and stores it in a variable declared as `Blacksmith` — the interface type, not `OrcBlacksmith`.
2. It calls `blacksmith.manufactureWeapon(WeaponType.SPEAR)`. Because the actual object behind that interface reference is an `OrcBlacksmith`, this runs `OrcBlacksmith`'s implementation, which looks up `SPEAR` in its pre-built map and returns an `OrcWeapon`.
3. The returned value is stored in a variable declared as `Weapon` — again the interface, not `OrcWeapon`. The calling code never needed to know or write `OrcWeapon` anywhere.
4. `App.main()` then constructs `new ElfBlacksmith()` and assigns it to the *same* `blacksmith` variable. From this point, calling `manufactureWeapon` on that same variable name now runs `ElfBlacksmith`'s implementation instead — same call site, different concrete behavior, because the object behind the interface changed.
5. This is the entire demonstration: one line changed (which concrete class got constructed), and every subsequent call using the interface type automatically produced different concrete results.

### Method-by-method: who calls what

| File | Method | Called by | It calls |
|---|---|---|---|
| `ElfBlacksmith.java` | `manufactureWeapon(type)` | `App.main()`, via the `Blacksmith` interface | nothing — reads from its pre-built `ELFARSENAL` map |
| `OrcBlacksmith.java` | `manufactureWeapon(type)` | `App.main()`, via the `Blacksmith` interface | nothing — reads from its pre-built `ORCARSENAL` map |
| `App.java` | `main()` | program entry point | `new OrcBlacksmith()` / `new ElfBlacksmith()`, then `manufactureWeapon(...)` |

---

## Go version

### The files and what each one is

| File | What it is |
|---|---|
| `weapon.go` | `Weapon` interface, `WeaponType`, `ElfWeapon`/`OrcWeapon` structs. |
| `blacksmith.go` | `Blacksmith` interface, `ElfBlacksmith`/`OrcBlacksmith` structs. |
| `main.go` | Demonstrates swapping which factory is held. |

Go's interfaces are **satisfied implicitly** — a struct doesn't declare
"I implement `Blacksmith`" anywhere; it satisfies the interface automatically
the moment it has methods matching the interface's signature. There's no
`implements` keyword at all.

### Who calls whom

```
main()
 └── holds a Blacksmith variable (static type is the INTERFACE)
        ├── could hold: ElfBlacksmith  ──produces──>  ElfWeapon
        └── could hold: OrcBlacksmith  ──produces──>  OrcWeapon
```

### What actually happens, in order, when you run it

1. `main()` assigns `NewOrcBlacksmith()` to a variable declared as `Blacksmith`.
2. Calling `.ManufactureWeapon(Spear)` on it runs `OrcBlacksmith`'s method (Go picks the right method automatically based on the concrete value stored inside the interface, called a "type switch" under the hood, though you never write one yourself here), returning an `OrcWeapon`.
3. `main()` reassigns the same variable to `NewElfBlacksmith()`. The next call to `.ManufactureWeapon(...)` on that variable now runs `ElfBlacksmith`'s method instead.
4. Same demonstration as Java: one line changes which concrete value is inside the interface variable, and every call through that variable changes behavior with it.

### Method-by-method: who calls what

| File | Method | Called by | It calls |
|---|---|---|---|
| `blacksmith.go` | `(ElfBlacksmith).ManufactureWeapon(type)` | `main()`, via the `Blacksmith` interface | nothing — reads from its own `arsenal` map |
| `blacksmith.go` | `(OrcBlacksmith).ManufactureWeapon(type)` | `main()`, via the `Blacksmith` interface | nothing — reads from its own `arsenal` map |
| `main.go` | `main()` | program entry point | `NewOrcBlacksmith()` / `NewElfBlacksmith()`, then `.ManufactureWeapon(...)` |

---

## Side-by-side: same idea, different word

| What it does | Java | Go |
|---|---|---|
| Declaring "this satisfies the interface" | explicit — `implements Blacksmith` | implicit — no keyword; matching method signatures is enough |
| The product interface | `Weapon` | `Weapon` |
| The creator interface | `Blacksmith` | `Blacksmith` |
| Choosing pointer vs. value for the concrete factories | not applicable (Java references are always reference semantics) | value receivers, deliberately — nothing here ever mutates after construction, so no pointer is needed |
| What actually decides which concrete product comes back | which subclass you instantiated | which concrete type is stored inside the interface variable |

---

## Pseudocode you can write from scratch in an interview

```
define Product interface: { describeYourself() }

define ConcreteProductA implements Product
define ConcreteProductB implements Product

define Creator interface: { makeProduct() -> Product }

define ConcreteCreatorA implements Creator:
    makeProduct() -> returns a ConcreteProductA

define ConcreteCreatorB implements Creator:
    makeProduct() -> returns a ConcreteProductB

# usage
creator: Creator = ConcreteCreatorA()
product = creator.makeProduct()   # caller never named ConcreteProductA directly
```

The one line worth memorizing: **"the caller only ever holds and calls
through the interface types — the moment you see a concrete class name
show up outside of a single construction line, the pattern's benefit is
being bypassed."**

---

## Likely interview questions, answered short

See `QA.md` in this same folder for the full interview-style question set
with model answers.
