package models

type UserBattleHistory struct {
	UserID int
	NumberOfBattles int
	BattlesWon int
	BattleLost int
	Trophies int
}