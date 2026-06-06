package models

type Troop struct {
	TroopID int
	Name string
	BaseCost int
	TroopAttackPower int
	BuildingAttackPower int
	Defence int
	Range int
	AttributeStrength int
	AttributeWeakness int
	TroopArmySpace int
	MovementSpeed int
	MaxLevel []int
}