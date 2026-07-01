package controller

import (
	"encoding/json"
	"net/http"
	"rats/models"
    "log"
)

func GetVillage(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)
    models.UpdatePassiveResources(userID) 

    buildings, err := models.GetVillageBuildings(userID)
    if err != nil {
        log.Println("GetVillage error:", err)
        http.Error(w, "Failed to fetch village", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(buildings)
}

func PlaceBuilding(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)
	var input struct{
		BuildingID int `json:"building_id"`
		X int `json:"x"`
		Y int `json:"y"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	err := models.PlaceBuilding(userID, input.BuildingID, input.X , input.Y)
	if err != nil {
        http.Error(w, "Failed to Place Building", http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "building placed"})

}

func MoveBuilding(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)

    var input struct {
        InstanceID int `json:"instance_id"`
        X          int `json:"x"`
        Y          int `json:"y"`
    }
    json.NewDecoder(r.Body).Decode(&input)

    // verify this building belongs to this user
    belongs, err := models.BuildingBelongsToUser(input.InstanceID, userID)
    if err != nil || !belongs {
        http.Error(w, "Building not found", http.StatusNotFound)
        return
    }

    err = models.MoveBuilding(input.InstanceID, input.X, input.Y)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "building moved"})
}

func UpgradeBuilding(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)

    var input struct {
        InstanceID int `json:"instance_id"`
    }
    json.NewDecoder(r.Body).Decode(&input)

    belongs, err := models.BuildingBelongsToUser(input.InstanceID, userID)
    if err != nil || !belongs {
        http.Error(w, "Building not found", http.StatusNotFound)
        return
    }

    err = models.UpgradeBuilding(input.InstanceID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "building upgraded"})
}

func GetCity(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)
    models.UpdatePassiveResources(userID)   
    city, err := models.GetCity(userID)
    if err != nil {
        log.Println("GetCity error:", err)
        http.Error(w, "Failed to fetch city", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(city)
}


func GetUserBattleHistory(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)
    battleHistory, err := models.GetUserBattleHistory(userID)
    if err != nil {
        log.Println("GetBattleHistory error:", err)
        http.Error(w, "Failed to fetch BattleHistory", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(battleHistory)
}