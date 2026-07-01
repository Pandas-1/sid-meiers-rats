package models

import (
	"rats/db"
)

type UserBattleHistory struct {
	UserID int `json:"user_id"`
	NumberOfBattles int `json:"number_of_battles"`
	BattlesWon int `json:"battles_won"`
	BattlesLost int `json:"battles_lost"`
	Trophies int `json:"trophies"`
}

func GetUserBattleHistory(userID int) (UserBattleHistory, error) {

	var battleHistory UserBattleHistory
	battleHistory.UserID = userID
	err := db.DB.QueryRow(
		"SELECT number_of_battles, battles_won, battles_lost, trophies FROM user_battle_history WHERE user_id = $1",
		userID,
	).Scan(&battleHistory.NumberOfBattles, &battleHistory.BattlesWon, &battleHistory.BattlesLost, &battleHistory.Trophies)
	return battleHistory, err

}
