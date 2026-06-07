package models

import (
	"time"
	"rats/db"
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
	_, err := db.DB.Exec(
        "INSERT INTO users (username, password_hash, trophies, base_level, xp, last_played) VALUES ($1, $2, 0, 1, 0, NOW())",
        username, passwordHash,
    )
    return err
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