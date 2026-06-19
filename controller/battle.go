package controller

import (
    "encoding/json"
    "net/http"
    "log"
    "rats/battle"
)

func StartBattle(w http.ResponseWriter, r *http.Request) {
    attackerID := r.Context().Value(UserIDKey).(int)

    var input struct {
        DefenderID int `json:"defender_id"`
    }
    json.NewDecoder(r.Body).Decode(&input)
     log.Printf("StartBattle: attackerID=%d defenderID=%d", attackerID, input.DefenderID)

    battleID, err := battle.StartBattle(attackerID, input.DefenderID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"battle_id": battleID})
}

func BattleWS(w http.ResponseWriter, r *http.Request) {
    battle.BattleWSHandler(w, r)
}