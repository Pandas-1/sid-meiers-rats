package controller

import (
    "encoding/json"
    "net/http"
    "log"
    "rats/battle"
    "rats/models"
)

func StartBattle(w http.ResponseWriter, r *http.Request) {
    attackerID := r.Context().Value(UserIDKey).(int)

    //var input struct {
    //    DefenderID int `json:"defender_id"`
    //}
    //json.NewDecoder(r.Body).Decode(&input)

    models.PendingMatchesMu.Lock()
    defenderID, ok := models.PendingMatches[attackerID]

    if ok {
        delete(models.PendingMatches, attackerID)
    }
    models.PendingMatchesMu.Unlock()

    if !ok {
        http.Error(w, "You must matchmake before attacking", 400)
        return
    }

    log.Printf("StartBattle: attackerID=%d defenderID=%d", attackerID, defenderID)

    battleID, err := battle.StartBattle(attackerID, defenderID)
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