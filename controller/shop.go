package controller

import (
    "encoding/json"
    "net/http"
    "rats/models"
)

func GetBuildings(w http.ResponseWriter, r *http.Request) {
    buildings, err := models.GetAllBuildingDetails()
    if err != nil {
        http.Error(w, "Failed to fetch buildings", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(buildings)
}

func GetTroops(w http.ResponseWriter, r *http.Request) {
    troops, err := models.GetAllTroopDetails()
    if err != nil {
        http.Error(w, "Failed to fetch troops", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(troops)
}