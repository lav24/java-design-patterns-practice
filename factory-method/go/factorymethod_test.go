package main

import "testing"

// verifyWeapon checks that weapon is the expected concrete type and carries
// the expected WeaponType — mirrors FactoryMethodTest.verifyWeapon in Java.
func verifyWeapon(t *testing.T, weapon Weapon, expectedType WeaponType, wantElf bool) {
	t.Helper()

	_, isElf := weapon.(ElfWeapon)
	_, isOrc := weapon.(OrcWeapon)
	if wantElf && !isElf {
		t.Errorf("weapon is not an ElfWeapon: %#v", weapon)
	}
	if !wantElf && !isOrc {
		t.Errorf("weapon is not an OrcWeapon: %#v", weapon)
	}
	if weapon.WeaponType() != expectedType {
		t.Errorf("weapon type = %v, want %v", weapon.WeaponType(), expectedType)
	}
}

func TestOrcBlacksmithWithSpear(t *testing.T) {
	blacksmith := NewOrcBlacksmith()
	weapon := blacksmith.ManufactureWeapon(Spear)
	verifyWeapon(t, weapon, Spear, false)
}

func TestOrcBlacksmithWithAxe(t *testing.T) {
	blacksmith := NewOrcBlacksmith()
	weapon := blacksmith.ManufactureWeapon(Axe)
	verifyWeapon(t, weapon, Axe, false)
}

func TestElfBlacksmithWithShortSword(t *testing.T) {
	blacksmith := NewElfBlacksmith()
	weapon := blacksmith.ManufactureWeapon(ShortSword)
	verifyWeapon(t, weapon, ShortSword, true)
}

func TestElfBlacksmithWithSpear(t *testing.T) {
	blacksmith := NewElfBlacksmith()
	weapon := blacksmith.ManufactureWeapon(Spear)
	verifyWeapon(t, weapon, Spear, true)
}
