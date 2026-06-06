package db

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

var DB *sql.DB 

func Connect() {
	conn, err := sql.Open("postgres", "host=localhost user=postgres password=Gedelaisstupid dbname=rats sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	
	err = conn.Ping()  
    if err != nil {
        log.Fatal(err)
    }

	DB = conn
}