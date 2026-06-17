package models

import(
	"rats/db"
	"fmt"
    "github.com/lib/pq"
)

type UserBuilding struct {
	InstanceID int
	UserID int
	BuildingID int
	Level int
	GridX int
	GridY int
    // gotta do this to make only one call for db to make the buildings
    Width        int
    Height       int
    Name         string
    BuildingType string
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
    rows, err := db.DB.Query(`
        SELECT ub.instance_id, ub.user_id, ub.building_id, ub.level, 
               ub.grid_x, ub.grid_y, bd.width, bd.height, bd.name, bd.building_type
        FROM user_buildings ub
        JOIN building_details bd ON ub.building_id = bd.building_id
        WHERE ub.user_id = $1
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var buildings []UserBuilding
    for rows.Next() {
        var b UserBuilding
        err := rows.Scan(&b.InstanceID, &b.UserID, &b.BuildingID, &b.Level, 
                        &b.GridX, &b.GridY, &b.Width, &b.Height, &b.Name, &b.BuildingType)
        if err != nil {
            return nil, err
        }
        buildings = append(buildings, b)
    }
    return buildings, nil
}

func MoveBuilding(InstanceID int, x int, y int) error {
    var width, height, userID int
    err := db.DB.QueryRow(
        `SELECT bd.width, bd.height, ub.user_id 
        FROM user_buildings ub
        JOIN building_details bd ON ub.building_id = bd.building_id
        WHERE ub.instance_id = $1`,
        InstanceID,
    ).Scan(&width, &height, &userID)
    // i gotta remove the orignal building when doing collision
    if err != nil {
        return err
    }
    var count int
    err = db.DB.QueryRow(
        `SELECT COUNT(*) FROM user_buildings ub
        JOIN building_details bd ON ub.building_id = bd.building_id
        WHERE ub.user_id = $1
        AND ub.grid_x < $2 AND ub.grid_x + bd.width > $3
        AND ub.grid_y < $4 AND ub.grid_y + bd.height > $5
        AND ub.instance_id != $6`,
        userID, x+width, x, y+height, y, InstanceID,
    ).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        return fmt.Errorf("grid space already occupied")
    }

    _, err = db.DB.Exec(
        "UPDATE user_buildings SET grid_x = $1 , grid_y = $2 WHERE instance_id = $3 ",
        x, y, InstanceID,
    )
    return err
}

// define func upgrade building later when i understand how i want to do the scaling
// maybe will use smth here for testing

func UpgradeBuilding(instanceID int) error {
    var currentLevel int
    var userID int
    var maxLevels pq.Int64Array
    var costResource1, costResource2, baseLevel , scaling int

    err := db.DB.QueryRow(
        `SELECT ub.level, ub.user_id, bd.max_level, bd.cost_resource1, bd.cost_resource2, users.base_level, bd.scaling
         FROM user_buildings ub
         JOIN building_details bd ON ub.building_id = bd.building_id
         JOIN users ON users.user_id = ub.user_id
         WHERE ub.instance_id = $1`,
        instanceID,
    ).Scan(&currentLevel, &userID, &maxLevels, &costResource1, &costResource2, &baseLevel, &scaling)
    if err != nil {
        return fmt.Errorf("building not found: %w", err)
    }

    // check if already at max level
    absoluteMax := int(maxLevels[baseLevel - 1])
    if currentLevel >= absoluteMax {
        return fmt.Errorf("building already at max level")
    }

    // check if user has enough resources
    var resource1, resource2 int64
    err = db.DB.QueryRow(
        "SELECT resource1, resource2 FROM city_details WHERE user_id = $1",
        userID,
    ).Scan(&resource1, &resource2)
    if err != nil {
        return fmt.Errorf("could not fetch resources: %w", err)
    }
    
    upgradeCostResource1 := costResource1*scaling*currentLevel
    upgradeCostResource2 := costResource2*scaling*currentLevel

    if int64(upgradeCostResource1) > resource1 || int64(upgradeCostResource2) > resource2 {
        return fmt.Errorf("not enough resources")
    }

    // deduct resources
    _, err = db.DB.Exec(
        "UPDATE city_details SET resource1 = resource1 - $2, resource2 = resource2 - $3 WHERE user_id = $1",
        userID, upgradeCostResource1, upgradeCostResource2,
    )
    if err != nil {
        return err
    }

    // increment level
    _, err = db.DB.Exec(
        "UPDATE user_buildings SET level = level + 1 WHERE instance_id = $1",
        instanceID,
    )
    return err
}

func BuildingBelongsToUser(instanceID, userID int) (bool, error) {
    var count int
    err := db.DB.QueryRow(
        "SELECT COUNT(*) FROM user_buildings WHERE instance_id = $1 AND user_id = $2",
        instanceID, userID,
    ).Scan(&count)
    return count > 0, err
}