package main

// Blacksmith is the creator interface — it declares the factory method.
// Callers only ever hold a Blacksmith and get back a Weapon; which
// concrete weapon type that is depends entirely on which concrete
// Blacksmith they're holding, not on anything the caller does.
type Blacksmith interface {
	ManufactureWeapon(weaponType WeaponType) Weapon
	String() string
}

// ElfBlacksmith always returns ElfWeapon — the factory method's decision
// is fixed by which struct this is, not by anything passed into it.
// Value receivers throughout: nothing here is ever mutated after
// construction, so there's no reason to use a pointer.
type ElfBlacksmith struct {
	arsenal map[WeaponType]Weapon
}

func NewElfBlacksmith() ElfBlacksmith {
	arsenal := make(map[WeaponType]Weapon)
	for _, t := range []WeaponType{ShortSword, Spear, Axe} {
		arsenal[t] = ElfWeapon{weaponType: t}
	}
	return ElfBlacksmith{arsenal: arsenal}
}

func (b ElfBlacksmith) ManufactureWeapon(weaponType WeaponType) Weapon {
	return b.arsenal[weaponType]
}

func (b ElfBlacksmith) String() string { return "The elf blacksmith" }

// OrcBlacksmith always returns OrcWeapon.
type OrcBlacksmith struct {
	arsenal map[WeaponType]Weapon
}

func NewOrcBlacksmith() OrcBlacksmith {
	arsenal := make(map[WeaponType]Weapon)
	for _, t := range []WeaponType{ShortSword, Spear, Axe} {
		arsenal[t] = OrcWeapon{weaponType: t}
	}
	return OrcBlacksmith{arsenal: arsenal}
}

func (b OrcBlacksmith) ManufactureWeapon(weaponType WeaponType) Weapon {
	return b.arsenal[weaponType]
}

func (b OrcBlacksmith) String() string { return "The orc blacksmith" }
