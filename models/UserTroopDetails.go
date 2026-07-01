package models

import(
	"rats/db"
	"fmt"
	"github.com/lib/pq"
)

type UserTroopDetail struct {
	UserID int 
	TroopID int
	TroopLevel int
}

func GetTroopLevel(userID int, troopID int) (int, error) {
	var troopLevel int 
	row := db.DB.QueryRow("SELECT troop_level FROM user_troop_details WHERE user_id = $1 AND troop_id = $2",
	userID, troopID)
	err := row.Scan(&troopLevel)
	return troopLevel, err
	
}

func LevelUpTroop(userID int, troopID int) error {
	var baseLevel int
	var currentLevel int
	var maxLevels pq.Int64Array
	var baseCost int
	var scaling int = 2

	err := db.DB.QueryRow(`
		SELECT u.base_level, utd.troop_level, td.max_level, td.base_cost
		FROM users u
		JOIN user_troop_details utd ON u.user_id = utd.user_id
		JOIN troop_details td ON td.troop_id = $2
		WHERE u.user_id = $1 AND utd.troop_id = $2
	`, userID, troopID).Scan(&baseLevel, &currentLevel, &maxLevels, &baseCost)

	if err != nil {
		return fmt.Errorf("troop not found or not unlocked: %w", err)
	}
	
	absoluteMax := int(maxLevels[baseLevel-1])
	if currentLevel >= absoluteMax {
		return fmt.Errorf("troop already at max level")
	}

	var resource1, resource2 int64
	err = db.DB.QueryRow(
		"SELECT resource1, resource2 FROM city_details WHERE user_id = $1",
		userID,
	).Scan(&resource1, &resource2)
	
	if err != nil {
		return fmt.Errorf("could not fetch resources: %w", err)
	}

	// Using resource2 for base cost
	upgradeCost := baseCost * currentLevel * scaling

	if resource2 < int64(upgradeCost) {
		return fmt.Errorf("not enough resources")
	}

	_, err = db.DB.Exec(
		"UPDATE city_details SET resource2 = resource2 - $2 WHERE user_id = $1",
		userID, upgradeCost,
	)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(
		"UPDATE user_troop_details SET troop_level = troop_level + 1 WHERE user_id = $1 AND troop_id = $2",
		userID, troopID,
	)
	
	return err
}


func GetUserTroops(userID int) ([]UserTroopDetail, error) {
    rows, err := db.DB.Query(
        "SELECT user_id, troop_id, troop_level FROM user_troop_details WHERE user_id = $1",
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var troops []UserTroopDetail
    for rows.Next() {
        var t UserTroopDetail
        err := rows.Scan(&t.UserID, &t.TroopID, &t.TroopLevel)
        if err != nil {
            return nil, err
        }
        troops = append(troops, t)
    }
    return troops, nil
}