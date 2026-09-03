package main

// WeaponType mirrors the Java enum — a closed set of weapon kinds.
type WeaponType int

const (
	ShortSword WeaponType = iota
	Spear
	Axe
)

func (w WeaponType) String() string {
	switch w {
	case ShortSword:
		return "short sword"
	case Spear:
		return "spear"
	case Axe:
		return "axe"
	default:
		return ""
	}
}

// Weapon is the product interface — the only thing callers ever see,
// regardless of which concrete factory built it.
type Weapon interface {
	WeaponType() WeaponType
}

// ElfWeapon and OrcWeapon are the two concrete products. Structurally
// identical — they only differ in how they describe themselves — same as
// the Java ElfWeapon/OrcWeapon records.
type ElfWeapon struct {
	weaponType WeaponType
}

func (w ElfWeapon) WeaponType() WeaponType { return w.weaponType }
func (w ElfWeapon) String() string         { return "an elven " + w.weaponType.String() }

type OrcWeapon struct {
	weaponType WeaponType
}

func (w OrcWeapon) WeaponType() WeaponType { return w.weaponType }
func (w OrcWeapon) String() string         { return "an orcish " + w.weaponType.String() }
