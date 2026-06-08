package models

import(
	"rats/db"
	"fmt"
)

type UserBuilding struct {
	InstanceID int
	UserID int
	BuildingID int
	Level int
	GridX int
	GridY int
}

func PlaceBuilding(userID, buildingID, x, y int) error {
    //  get building size
    var width, height int
    err := db.DB.QueryRow(
        "SELECT width, height FROM building_details WHERE building_id = $1",
        buildingID,
    ).Scan(&width, &height)
    if err != nil {
        return fmt.Errorf("building not found: %w", err)
    }

    //  overlap
    var count int
    err = db.DB.QueryRow(
        `SELECT COUNT(*) FROM user_buildings ub
         JOIN building_details bd ON ub.building_id = bd.building_id
         WHERE ub.user_id = $1
         AND ub.grid_x < $2 AND ub.grid_x + bd.width > $3
         AND ub.grid_y < $4 AND ub.grid_y + bd.height > $5`,
        userID, x+width, x, y+height, y,
    ).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        return fmt.Errorf("grid space already occupied")
    }

    // insert
    _, err = db.DB.Exec(
        "INSERT INTO user_buildings (user_id, building_id, level, grid_x, grid_y) VALUES ($1, $2, 1, $3, $4)",
        userID, buildingID, x, y,
    )
    return err
}

func GetVillageBuildings(userID int) ([]UserBuilding, error) {
    buildings := make([]UserBuilding, 0)

    rows, err := db.DB.Query(
        `SELECT instance_id, user_id, building_id, level, grid_x, grid_y
         FROM user_buildings
         WHERE user_id = $1`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var b UserBuilding

        err := rows.Scan(
            &b.InstanceID,
            &b.UserID,
            &b.BuildingID,
            &b.Level,
            &b.GridX,
            &b.GridY,
        )
        if err != nil {
            return nil, err
        }

        buildings = append(buildings, b)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return buildings, nil
}