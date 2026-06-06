package models

type Battle struct {
	BattleID int
	AttackID int
	DefenderId int
	Resource1Won int
	Resource2Won int
	VictoryPercentage int
}