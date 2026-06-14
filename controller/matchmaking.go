package controller

import (
    "encoding/json"
    "net/http"
    "rats/models"
)

func FindOpponent(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(int)

    opponent, err := models.FindOpponent(userID)
    if err != nil {
        http.Error(w, "No opponent found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(opponent)
}