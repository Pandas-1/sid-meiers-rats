package models

import (
	"time"
	"rats/db"
)

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


func CreateCity(userID int) error {
	_ , err := db.DB.Exec(
		"INSERT into city_details (user_id, resource1 , resource2, max_resource1 , max_resource2, max_troop_army_size, last_updated, max_defence_buildings, max_resource_buildings) VALUES ($1, 100 , 100 , 1000, 1000, 30, NOW(), 4, 4) ", 
		userID,
	)
	return err
} 

func GetCity(userID int) (CityDetails, error) {
	var c CityDetails
	c.UserID = userID
	row := db.DB.QueryRow(
		"SELECT resource1,resource2,max_resource1,max_resource2,max_troop_army_size,last_updated,max_defence_buildings, max_resource_buildings FROM city_details WHERE user_id = $1",
		userID,
	)
	err := row.Scan(&c.Resource1, &c.Resource2, &c.MaxResource1, &c.MaxResource2, &c.MaxTroopArmySize, &c.LastUpdated, &c.MaxDefenceBuildings, &c.MaxResourceBuildings)
    return c, err
}

func UpdateResources(userID int, resource1, resource2 int64) error {
	_ , err := db.DB.Exec(
		"UPDATE city_details SET resource1 = $2, resource2 = $3 WHERE user_id = $1",
		userID, resource1, resource2,
	)
	return err

}

func AddResources(userID int, resource1, resource2 int64) error {
    _, err := db.DB.Exec(
        `UPDATE city_details 
         SET resource1 = GREATEST(0, resource1 + $2),
             resource2 = GREATEST(0, resource2 + $3)
         WHERE user_id = $1`,
        userID, resource1, resource2,
    )
    return err
}