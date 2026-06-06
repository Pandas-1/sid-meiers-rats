package models

import "time"

type CityDetails struct {
	UserID int
	Resource1 int
	Resource2 int
	MaxResource1 int
	MaxResource2 int
	MaxTroopArmySize int
	LastUpdated time.Time
	MaxDefenceBuildings int 
	MaxResourceBuildings int
}