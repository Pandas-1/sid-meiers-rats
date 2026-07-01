package models

import (
	"rats/db"
	"time"
	"fmt"
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
	// REMIND YOURSELF TO DO THE TROPHIES PARAMETER LATER
	_, err := db.DB.Exec(
		`INSERT INTO battles (attacker_id, defender_id, resource1_won, resource2_won, victory_percentage) 
		 VALUES ($1, $2, $3, $4, $5)`,
		attackerID, defenderID, r1won, r2won, victoryPct,
	)
	if err != nil {
		return fmt.Errorf("failed to insert battle: %w", err)
	}

	var attackerWon, attackerLost, attackerTrophyChange int
	var defenderWon, defenderLost, defenderTrophyChange int

	if victoryPct > 50 {
		attackerWon, attackerLost = 1, 0
		defenderWon, defenderLost = 0, 1
		attackerTrophyChange = +25
		defenderTrophyChange = -15
	} else {
		attackerWon, attackerLost = 0, 1
		defenderWon, defenderLost = 1, 0
		attackerTrophyChange= -10
		defenderTrophyChange = 10
	}

	_, err = db.DB.Exec(`
		INSERT INTO user_battle_history (user_id, number_of_battles, battles_won, battles_lost, trophies)
		VALUES ($1, 1, $2, $3, 0)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			number_of_battles = user_battle_history.number_of_battles + 1,
			battles_won = user_battle_history.battles_won + $2,
			battles_lost = user_battle_history.battles_lost + $3,
			trophies = GREATEST(0, user_battle_history.trophies + $4)
	`, attackerID, attackerWon, attackerLost, attackerTrophyChange)

	if err != nil {
		return fmt.Errorf("failed to update attacker history: %w", err)
	}

	_, err = db.DB.Exec(`
		INSERT INTO user_battle_history (user_id, number_of_battles, battles_won, battles_lost, trophies)
		VALUES ($1, 1, $2, $3, 0)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			number_of_battles = user_battle_history.number_of_battles + 1,
			battles_won = user_battle_history.battles_won + $2,
			battles_lost = user_battle_history.battles_lost + $3,
			trophies = GREATEST(0, user_battle_history.trophies + $4)
	`, defenderID, defenderWon, defenderLost, defenderTrophyChange)

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