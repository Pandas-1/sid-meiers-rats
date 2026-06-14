package models

import(
	"rats/db"
	"github.com/lib/pq"

)

type BuildingDetails struct {
	BuildingID int
	Name string
	BuildingType string
	Production int
	Scaling int
	HealthBar int
	Width int
	Height int
	DefenceAttack int
	DefenceRange int
	MaxLevel pq.Int64Array
	CostResource1 int
	CostResource2 int

}

func GetBuildingDetails(buildingID int) (BuildingDetails, error) {
	var b BuildingDetails

	row := db.DB.QueryRow(
		`SELECT building_id, name, building_type, production, scaling,
		        health_bar, width, height, defence_attack,
		        defence_range, max_level, cost_resource1, cost_resource2
		   FROM building_details
		  WHERE building_id = $1`,
		buildingID,
	)

	err := row.Scan(
		&b.BuildingID,
		&b.Name,
		&b.BuildingType,
		&b.Production,
		&b.Scaling,
		&b.HealthBar,
		&b.Width,
		&b.Height,
		&b.DefenceAttack,
		&b.DefenceRange,
		&b.MaxLevel,
		&b.CostResource1,
		&b.CostResource2,
	)

	return b, err
}

func GetAllBuildingDetails() ([]BuildingDetails, error) {
    rows, err := db.DB.Query("SELECT building_id, name, building_type, production, scaling, health_bar, width, height, defence_attack, defence_range, max_level, cost_resource1, cost_resource2 FROM building_details")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var buildings []BuildingDetails
    for rows.Next() {
        var b BuildingDetails
        err := rows.Scan(&b.BuildingID, &b.Name, &b.BuildingType, &b.Production, &b.Scaling, &b.HealthBar, &b.Width, &b.Height, &b.DefenceAttack, &b.DefenceRange, &b.MaxLevel, &b.CostResource1, &b.CostResource2)
        if err != nil {
            return nil, err
        }
        buildings = append(buildings, b)
    }
    return buildings, nil
}