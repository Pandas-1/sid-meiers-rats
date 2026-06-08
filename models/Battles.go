package models

import (
	"rats/db"
	"time"
)

type Battle struct {
	BattleID          int
	AttackerID        int
	DefenderID        int
	Resource1Won      int
	Resource2Won      int
	VictoryPercentage int
	FoughtAt          time.Time
}

func CreateBattle(attackerID, defenderID, r1won, r2won, victoryPct int) error {
	_, err := db.DB.Exec(
		`INSERT INTO battles (attacker_id, defender_id, resource1_won, resource2_won, victory_percentage) 
		 VALUES ($1, $2, $3, $4, $5)`,
		attackerID, defenderID, r1won, r2won, victoryPct,
	)
	return err
}

func GetBattleHistory(userID int) ([]Battle, error) {
	rows, err := db.DB.Query(
		`SELECT battle_id, attacker_id, defender_id, resource1_won, resource2_won, victory_percentage, fought_at 
		 FROM battles 
		 WHERE attacker_id = $1 OR defender_id = $1 
		 ORDER BY fought_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close() 

	var history []Battle

	for rows.Next() {
		var b Battle
		err := rows.Scan(
			&b.BattleID,
			&b.AttackerID,
			&b.DefenderID,
			&b.Resource1Won,
			&b.Resource2Won,
			&b.VictoryPercentage,
			&b.FoughtAt,
		)
		if err != nil {
			return nil, err
		}
		// Add the single battle we just scanned into our history slice
		history = append(history, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}