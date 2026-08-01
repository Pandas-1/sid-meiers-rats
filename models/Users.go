package models

import (
	"time"
	"rats/db"
	"fmt"
)

type User struct {
	UserID int
	Username string
	Trophies int
	BaseLevel int
	XP int
	LastPlayed time.Time
	PasswordHash string
}

func CreateUser(username, passwordHash string) error {
    // start transaction
    tx, err := db.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    var userID int
    err = tx.QueryRow(
        "INSERT INTO users (username, password_hash, base_level, xp, last_played) VALUES ($1, $2, 1, 0, NOW()) RETURNING user_id",
        username, passwordHash,
    ).Scan(&userID)
    if err != nil {
        return fmt.Errorf("could not create user: %w", err)
    }

    //create city_details entry, bcz i am fucking with transactions idk if i can call the CreateCity but i dont really care i think
    _, err = tx.Exec(
        "INSERT INTO city_details (user_id, resource1, resource2, max_resource1, max_resource2, max_troop_army_size, last_updated, max_defence_buildings, max_resource_buildings) VALUES ($1, 100, 100, 1000, 1000, 30, NOW(), 4, 4)",
        userID,
    )
    if err != nil {
        return fmt.Errorf("could not create city: %w", err)
    }

    //create user_battle_history entry
    _, err = tx.Exec(
        "INSERT INTO user_battle_history (user_id, number_of_battles, battles_won, battles_lost, trophies) VALUES ($1, 0, 0, 0, 0)",
        userID,
    )
    if err != nil {
        return fmt.Errorf("could not create battle history: %w", err)
    }

    //reate army_details entry
    _, err = tx.Exec(
        "INSERT INTO army_details (user_id, troop_units_used, army_composition, created_on) VALUES ($1, 0, '[]', NOW())",
        userID,
    )
    if err != nil {
        return fmt.Errorf("could not create army: %w", err)
    }

    // create user_troop_details for ALL troops that exist in troop_details
    _, err = tx.Exec(
        `INSERT INTO user_troop_details (user_id, troop_id, troop_level)
         SELECT $1, troop_id, 1 FROM troop_details`,
        userID,
    )

    // forgor i also gotta add the townhall at the centre 2nd node, also gotta add elixer and mines
    _, err = tx.Exec(`
        INSERT INTO user_buildings (user_id, building_id, level, grid_x, grid_y)
        VALUES 
            ($1, (SELECT building_id FROM building_details WHERE name = 'Town Hall'), 1, 23, 23),
            ($1, (SELECT building_id FROM building_details WHERE name = 'Gold Mine'), 1, 10, 20),
            ($1, (SELECT building_id FROM building_details WHERE name = 'Elixir Collector'), 1, 20, 20)
    `, userID)
    if err != nil {
        return fmt.Errorf("could not place starting bildings: %w", err)
    }

    // all steps succeeded, commit
    return tx.Commit()
}

func GetUserByUsername(username string) (User, error) {
    var u User
    row := db.DB.QueryRow(
        "SELECT user_id, username, trophies, base_level, xp, last_played, password_hash FROM users WHERE username = $1",
        username,
    )
    err := row.Scan(&u.UserID, &u.Username, &u.Trophies, &u.BaseLevel, &u.XP, &u.LastPlayed, &u.PasswordHash)
    return u, err
}

func UpdateLastPlayed(userID int) error {
	_ , err := db.DB.Exec(
		"UPDATE users SET last_played = NOW() WHERE user_id = $1",
		userID,
	)
	return err
}