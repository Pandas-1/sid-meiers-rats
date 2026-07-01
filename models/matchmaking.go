package models

import (
	"rats/db"
	"fmt"
    "sync"
)

type MatchmakingResult struct {
    OpponentID       int    `json:"opponent_id"`
    OpponentUsername string `json:"opponent_username"`
    OpponentTrophies int    `json:"opponent_trophies"`
}

var PendingMatches = map[int]int{} // attackerID -> defenderID
var PendingMatchesMu sync.Mutex

func FindOpponent(userID int) (MatchmakingResult, error) {
    var userTrophies int
    err := db.DB.QueryRow(
        "SELECT trophies FROM user_battle_history WHERE user_id = $1",
        userID,
    ).Scan(&userTrophies)
    if err != nil {
        return MatchmakingResult{}, err
    }

    var result MatchmakingResult

    // try tight range first ±100
    err = db.DB.QueryRow(`
        SELECT u.user_id, u.username, ubh.trophies
        FROM users u
        JOIN user_battle_history ubh ON u.user_id = ubh.user_id
        WHERE u.user_id != $1
        AND ubh.trophies BETWEEN $2 AND $3
        ORDER BY RANDOM()
        LIMIT 1
    `, userID, userTrophies-100, userTrophies+100,
    ).Scan(&result.OpponentID, &result.OpponentUsername, &result.OpponentTrophies)

    if err != nil {
    // widen to ±300 if nobody found
    err = db.DB.QueryRow(`
        SELECT u.user_id, u.username, ubh.trophies
        FROM users u
        JOIN user_battle_history ubh ON u.user_id = ubh.user_id
        WHERE u.user_id != $1
        AND ubh.trophies BETWEEN $2 AND $3
        ORDER BY RANDOM()
        LIMIT 1
    `, userID, userTrophies-300, userTrophies+300,
    ).Scan(&result.OpponentID, &result.OpponentUsername, &result.OpponentTrophies)

    if err != nil {
        return MatchmakingResult{}, fmt.Errorf("no opponent found")
    }
    }
    PendingMatchesMu.Lock()
    PendingMatches[userID] = result.OpponentID
    PendingMatchesMu.Unlock()

    return result, nil
}