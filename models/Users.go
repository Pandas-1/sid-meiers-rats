package models

import "time"

type User struct {
	UserID int
	Username string
	Trophies int
	BaseLevel int
	Xp int
	LastPlayed time.Time
	PasswordHash string
}