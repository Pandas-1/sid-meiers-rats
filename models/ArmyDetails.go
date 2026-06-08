package models

import (
	"encoding/json"
	"fmt"
	"rats/db"
	"time"
)

type ArmyComposition struct {
	TroopID  int `json:"troop_id"`
	Quantity int `json:"quantity"`
}

type Army struct {
	UserID          int
	TroopUnitsUsed  int
	ArmyComposition []ArmyComposition
	CreatedOn       time.Time
}

func CreateOrUpdateArmy(userID int, composition []ArmyComposition) error {
	var totalUnits int

	for _, comp := range composition {
		var spacePerUnit int
		err := db.DB.QueryRow(
			"SELECT troop_army_space FROM troop_details WHERE troop_id = $1",
			comp.TroopID,
		).Scan(&spacePerUnit)
		
		if err != nil {
			return fmt.Errorf("invalid troop id %d: %w", comp.TroopID, err)
		}
		totalUnits += spacePerUnit * comp.Quantity
	}

	var maxCapacity int
	err := db.DB.QueryRow(
		"SELECT max_troop_army_size FROM city_details WHERE user_id = $1", 
		userID,
	).Scan(&maxCapacity)
	
	if err != nil {
		return fmt.Errorf("could not fetch city capacity: %w", err)
	}
	
	if totalUnits > maxCapacity {
		return fmt.Errorf("army size (%d) exceeds max capacity (%d)", totalUnits, maxCapacity)
	}

	// Convert the Go slice into JSON format
	jsonData, err := json.Marshal(composition)
	if err != nil {
		return fmt.Errorf("failed to process army composition: %w", err)
	}

	// Upsert the army into the database
	_, err = db.DB.Exec(`
		INSERT INTO army_details (user_id, troop_units_used, army_composition, created_on)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			troop_units_used = $2,
			army_composition = $3,
			created_on = NOW()
	`, userID, totalUnits, jsonData)

	return err
}

func GetArmy(userID int) (Army, error) {
	var a Army
	var rawJSON []byte

	err := db.DB.QueryRow(
		"SELECT user_id, troop_units_used, army_composition, created_on FROM army_details WHERE user_id = $1",
		userID,
	).Scan(&a.UserID, &a.TroopUnitsUsed, &rawJSON, &a.CreatedOn)
	if err != nil {
		return a, fmt.Errorf("army not found: %w", err)
	}

	err = json.Unmarshal(rawJSON, &a.ArmyComposition)
	if err != nil {
		return a, fmt.Errorf("failed to parse army composition: %w", err)
	}

	return a, nil
}

func ClearArmy(userID int) error {
	_, err := db.DB.Exec(
		`UPDATE army_details 
		 SET army_composition = '[]', troop_units_used = 0, created_on = NOW()
		 WHERE user_id = $1`,
		userID,
	)
	return err
}