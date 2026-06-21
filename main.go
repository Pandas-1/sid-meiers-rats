package main

import (
    "fmt"
    "net/http"
    "rats/db"
    "rats/routes"
)

func main() {
    db.Connect()
    fmt.Println("Connected to database!")

    r := routes.SetupRoutes()
    fmt.Println("Server starting on port 8080...")
    http.ListenAndServe(":8080", r)
}