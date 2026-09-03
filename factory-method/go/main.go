package main

import "fmt"

func main() {
	var blacksmith Blacksmith = NewOrcBlacksmith()
	weapon := blacksmith.ManufactureWeapon(Spear)
	fmt.Printf("%s manufactured %s\n", blacksmith, weapon)
	weapon = blacksmith.ManufactureWeapon(Axe)
	fmt.Printf("%s manufactured %s\n", blacksmith, weapon)

	blacksmith = NewElfBlacksmith()
	weapon = blacksmith.ManufactureWeapon(Spear)
	fmt.Printf("%s manufactured %s\n", blacksmith, weapon)
	weapon = blacksmith.ManufactureWeapon(Axe)
	fmt.Printf("%s manufactured %s\n", blacksmith, weapon)
}
