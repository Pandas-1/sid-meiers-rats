package db

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"os"
    "fmt"
)

var DB *sql.DB 

func Connect() {
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "localhost"  
    }
    user := os.Getenv("DB_USER")
    if user == "" {
        user = "postgres"
    }
    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        password = "Gedelaisstupid"
    }
    dbname := os.Getenv("DB_NAME")
    if dbname == "" {
        dbname = "rats"
    }

    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        host, user, password, dbname)

    conn, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = conn.Ping()
    if err != nil {
        log.Fatal(err)
    }

    DB = conn
    log.Println("Connected to database!")
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