package main

import(
	"fmt"
	"rats/db"
)

func main() {
	db.Connect()
	fmt.Println("Connected to database!")
}