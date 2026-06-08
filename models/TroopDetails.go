package models

import (
    "rats/db"
    "github.com/lib/pq"
)

type Troop struct {
    TroopID             int
    Name                string
    BaseCost            int
    TroopAttackPower    int
    BuildingAttackPower int
    Defence             int
    Range               int
    AttributeStrength   int
    AttributeWeakness   int
    TroopArmySpace      int
    MovementSpeed       int
    MaxLevel            pq.Int64Array
}


func GetTroopDetails(troopID int) (Troop, error) {
    var t Troop

    row := db.DB.QueryRow(
        `SELECT troop_id, name, base_cost, troop_attack_power, building_attack_power, 
                defence, range, attribute_strength, attribute_weakness, 
                troop_army_space, movement_speed, max_level
           FROM troop_details
          WHERE troop_id = $1`,
        troopID,
    )

    err := row.Scan(
        &t.TroopID,
        &t.Name,
        &t.BaseCost,
        &t.TroopAttackPower,
        &t.BuildingAttackPower,
        &t.Defence,
        &t.Range,
        &t.AttributeStrength,
        &t.AttributeWeakness,
        &t.TroopArmySpace,
        &t.MovementSpeed,
        &t.MaxLevel,
    )

    return t, err
}