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
	var totalCost int

	// start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// refund old army cost
	var oldCompositionRaw []byte
	err = tx.QueryRow(
		"SELECT army_composition FROM army_details WHERE user_id = $1",
		userID,
	).Scan(&oldCompositionRaw)

	// only refund if army already exists
	if err == nil {
		var oldComposition []ArmyComposition
		json.Unmarshal(oldCompositionRaw, &oldComposition)

		var refundAmount int
		for _, comp := range oldComposition {
			var baseCost, troopLevel int
			err := tx.QueryRow(
				`SELECT td.base_cost, utd.troop_level
				FROM troop_details td
				JOIN user_troop_details utd ON td.troop_id = utd.troop_id
				WHERE td.troop_id = $1 AND utd.user_id = $2`,
				comp.TroopID, userID,
			).Scan(&baseCost, &troopLevel)
			if err != nil {
				continue
			}
			refundAmount += baseCost * troopLevel * comp.Quantity
		}

		// refund to resource2
		_, err = tx.Exec(
			"UPDATE city_details SET resource2 = resource2 + $2 WHERE user_id = $1",
			userID, refundAmount,
		)
		if err != nil {
			return fmt.Errorf("could not refund resources: %w", err)
		}
	}

	for _, comp := range composition {
		var spacePerUnit, baseCost, troopLevel int
		if comp.Quantity < 0 {
            return fmt.Errorf("invalid quantity for troop %d", comp.TroopID)
        }
		err := tx.QueryRow(
			`SELECT td.troop_army_space, td.base_cost, utd.troop_level
			 FROM troop_details td
			 JOIN user_troop_details utd ON td.troop_id = utd.troop_id
			 WHERE td.troop_id = $1 AND utd.user_id = $2`,
			comp.TroopID, userID,
		).Scan(&spacePerUnit, &baseCost, &troopLevel)
		if err != nil {
			return fmt.Errorf("invalid troop id %d: %w", comp.TroopID, err)
		}

		totalUnits += spacePerUnit * comp.Quantity
		totalCost += baseCost * troopLevel * comp.Quantity
	}

	// check capacity
	var maxCapacity int
	err = tx.QueryRow(
		"SELECT max_troop_army_size FROM city_details WHERE user_id = $1",
		userID,
	).Scan(&maxCapacity)
	if err != nil {
		return fmt.Errorf("could not fetch city capacity: %w", err)
	}
	if totalUnits > maxCapacity {
		return fmt.Errorf("army size (%d) exceeds max capacity (%d)", totalUnits, maxCapacity)
	}

	// check resources
	var resource2 int64
	err = tx.QueryRow(
		"SELECT resource2 FROM city_details WHERE user_id = $1",
		userID,
	).Scan(&resource2)
	if err != nil {
		return fmt.Errorf("could not fetch resources: %w", err)
	}
	if int64(totalCost) > resource2 {
		return fmt.Errorf("not enough elixir — need %d have %d", totalCost, resource2)
	}

	// deduct resources
	_, err = tx.Exec(
		"UPDATE city_details SET resource2 = resource2 - $2 WHERE user_id = $1",
		userID, totalCost,
	)
	if err != nil {
		return err
	}

	// upsert army
	jsonData, err := json.Marshal(composition)
	if err != nil {
		return fmt.Errorf("failed to process army composition: %w", err)
	}
	_, err = tx.Exec(`
		INSERT INTO army_details (user_id, troop_units_used, army_composition, created_on)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			troop_units_used = $2,
			army_composition = $3,
			created_on = NOW()
	`, userID, totalUnits, jsonData)
	if err != nil {
		return err
	}

	return tx.Commit()
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