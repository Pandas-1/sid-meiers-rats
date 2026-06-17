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

func GameEntityData() {
    seedSQL, err := os.ReadFile("migrations/game_entity_data.sql")
    if err != nil {
        log.Fatal("Could not read seed file:", err)
    }

    _, err = DB.Exec(string(seedSQL))
    if err != nil {
        log.Fatal("Could not seed database:", err)
    }

    log.Println("Database seeded successfully!")
}