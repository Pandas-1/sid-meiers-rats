package controller

import (
	"net/http"
	"encoding/json"
	"rats/models"
)

func GetArmy(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)

    army, err := models.GetArmy(userID)
    if err != nil {
        http.Error(w, "Failed to fetch army", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(army)
}

func TrainArmy(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)

    // this is what client sends
    var input struct {
        Composition []models.ArmyComposition `json:"composition"`
    }
    json.NewDecoder(r.Body).Decode(&input)

    if len(input.Composition) == 0 {
        http.Error(w, "Army composition cannot be empty", http.StatusBadRequest)
        return
    }

    err := models.CreateOrUpdateArmy(userID, input.Composition)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "army trained"})
}

func LevelUpTroop(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)

    var input struct {
        TroopID int `json:"troop_id"`
    }
    json.NewDecoder(r.Body).Decode(&input)

    err := models.UpgradeTroop(userID, input.TroopID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "troop upgraded"})
}