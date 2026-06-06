package models

import (
	"time"
	"encoding/json"
)	

type ArmyComposition struct {
    TroopID  int `json:"troop_id"`
    Quantity int `json:"quantity"`
}


type Army struct {
	UserID int
	TroopUnitsUsed int
	ArmyComposition []ArmyComposition
	CreatedOn time.Time
} 