package battle

import (
    "encoding/json"
    "net/http"
    "sync"
    "rats/models"
    "github.com/gorilla/websocket"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "fmt"
    "log"
)

var jwtSecret = []byte("inthehistoryofjoeoverithasneverbeenthisjoeveridontevenknowifthisissecurebutafterenoughlettersithastobesecureatthispointlikewhat")

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins for now
    },
}

// store all active battle sessions
var (
    activeSessions   = map[int]*BattleSession{} // battleID → session
    sessionsMu       sync.Mutex
    nextBattleID     = 1
)

func StartBattle(attackerID, defenderID int) (int, error) {
    army, err := models.GetArmy(attackerID)
    if err != nil || len(army.ArmyComposition) == 0 {
        return 0, fmt.Errorf("you must train an army before attacking")
    }

    session, err := NewSession(attackerID, defenderID)
    if err != nil {
        return 0, err
    }
    log.Printf("Session created with %d buildings", session.InitBuildings)

    sessionsMu.Lock()
    battleID := nextBattleID
    nextBattleID++
    activeSessions[battleID] = session
    sessionsMu.Unlock()

    // start the game loop in background
    go runSession(battleID, session)

    return battleID, nil
}

func runSession(battleID int, session *BattleSession) {
    for i := 0; i < 200; i++ {
        if session.OnTick != nil {
            break
        }
        time.Sleep(50 * time.Millisecond)
    }
    for range session.Ticker.C {
        session.Tick()
        if session.OnTick != nil {
            session.OnTick(session.GetState())
        }
        if session.Done {
            session.Ticker.Stop()
            saveBattleResult(session)
            sessionsMu.Lock()
            delete(activeSessions, battleID)
            sessionsMu.Unlock()
            return
        }
    }
}

func BattleWSHandler(w http.ResponseWriter, r *http.Request) {
    // get battleID from query param and token verification

        log.Printf("WS request received")

    tokenString := r.URL.Query().Get("token")
    if tokenString == "" {
        http.Error(w, "No token provided", 401)
        return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        http.Error(w, "Invalid token", 401)
        return
    }


    battleIDStr := r.URL.Query().Get("battle_id")
    var battleID int
    json.Unmarshal([]byte(battleIDStr), &battleID)

    // get session
    sessionsMu.Lock()
    session, exists := activeSessions[battleID]
    sessionsMu.Unlock()

    if !exists {
        http.Error(w, "Battle not found", 404)
        return
    }

    // upgrade to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    session.ClientConnected = true 
    defer conn.Close()

    // set callback — every tick sends state to this client
    session.OnTick = func(state BattleState) {
        data, _ := json.Marshal(state)
        conn.WriteMessage(websocket.TextMessage, data)
    }

    // read troop drops from client
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            // client disconnected
            session.OnTick = nil
            break
        }

        var drop TroopDrop
        if err := json.Unmarshal(msg, &drop); err == nil {
            session.DropTroop(drop)
        }
    }
}

func saveBattleResult(session *BattleSession) {
    destroyed := session.InitBuildings - len(session.Buildings)
    destruction := 0
    if session.InitBuildings > 0 {
        destruction = (destroyed * 100) / session.InitBuildings
    }

    // calculate loot
    city, err := models.GetCity(session.DefenderID)
    if err != nil {
        return
    }
    loot1 := int(float64(city.Resource1) * 0.3)
    loot2 := int(float64(city.Resource2) * 0.3)

    models.CreateBattle(session.AttackerID, session.DefenderID, loot1, loot2, destruction)

    if destruction > 50 {
        models.AddResources(session.AttackerID, int64(loot1), int64(loot2))
        models.AddResources(session.DefenderID, int64(-loot1), int64(-loot2))
    }

    err = models.ClearArmy(session.AttackerID)
}