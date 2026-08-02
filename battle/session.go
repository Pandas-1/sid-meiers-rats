package battle

import (
	"math"
	"rats/models"
	"sync"
	"log"
	"time"
)

type BuildingState struct {
	InstanceID  int    `json:"instance_id"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	CurrentHP   int    `json:"current_hp"`
	MaxHP       int    `json:"max_hp"`
	Attack      int    `json:"attack"`
	Range       int    `json:"range"`
	Name        string `json:"name"`
	ElementType int    `json:"element_type"`
}

type TroopState struct {
	ID          int     `json:"id"`
	TroopID     int     `json:"troop_id"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	CurrentHP   int     `json:"current_hp"`
	MaxHP       int     `json:"max_hp"`
	Attack      int     `json:"attack"`
	Range       float64 `json:"range"`
	Speed       float64 `json:"speed"`
	Name        string  `json:"name"`
	ElementType int     `json:"element_type"`
}

type TroopDrop struct {
	TroopID int `json:"troop_id"`
	X       int `json:"x"`
	Y       int `json:"y"`
}

type BattleState struct {
	Buildings    []BuildingState `json:"buildings"`
	Troops       []TroopState    `json:"troops"`
	Done         bool            `json:"done"`
	Destruction  int             `json:"destruction"`
	ElapsedTicks int             `json:"elapsed_ticks"`
}

type BattleSession struct {
	AttackerID    int
	DefenderID    int
	Buildings     []BuildingState
	InitBuildings int
	Troops        []TroopState
	DropQueue     []TroopDrop
	Ticker        *time.Ticker
	Done          bool
	ElapsedTicks  int
	NextTroopID   int
	mu            sync.Mutex
	OnTick        func(BattleState)
	ArmyPool      map[int]models.Troop
	ArmyComp      []models.ArmyComposition
	ClientConnected bool
}

func NewSession(attackerID int, defenderID int) (*BattleSession, error) {
	villageBuildings, err := models.GetVillageBuildings(defenderID)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d buildings for defender %d", len(villageBuildings), defenderID)
    log.Printf("Buildings: %+v", villageBuildings)

	army, err := models.GetArmy(attackerID)
	if err != nil {
		return nil, err
	}

	armyPool := map[int]models.Troop{}
	for _, comp := range army.ArmyComposition {
		troopLevel, err := models.GetTroopLevel(attackerID, comp.TroopID)
		if err != nil {
			troopLevel = 1 // fallback, same as before
		}
		adjustedStats, err := models.GetTroopStatsAtLevel(comp.TroopID, troopLevel)
		if err != nil {
			continue
		}
		armyPool[comp.TroopID] = adjustedStats 
	}
	buildings := []BuildingState{}
	for _, b := range villageBuildings {
		adjustedStats, _ := models.GetBuildingStatsAtLevel(b.BuildingID, b.Level)
		buildings = append(buildings, BuildingState{
			InstanceID:  b.InstanceID,
			X:           b.GridX,
			Y:           b.GridY,
			Width:       b.Width,
			Height:      b.Height,
			CurrentHP:   adjustedStats.HealthBar,
			MaxHP:       adjustedStats.HealthBar,
			Attack:      adjustedStats.DefenceAttack,
			Range:       b.DefenceRange,
			Name:        b.Name,
			ElementType: b.ElementType,
		})
	}

	return &BattleSession{
		AttackerID:    attackerID,
		DefenderID:    defenderID,
		Buildings:     buildings,
		InitBuildings: len(buildings),
		Troops:        []TroopState{},
		DropQueue:     []TroopDrop{},
		Ticker:        time.NewTicker(time.Second / 20),
		Done:          false,
		ElapsedTicks:  0,
		NextTroopID:   1,
		ArmyPool:      armyPool,
		ArmyComp:      army.ArmyComposition,
	}, nil
}

func (s *BattleSession) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Done {
		return
	}

	s.ElapsedTicks++

	// place queued troop drops
	for _, drop := range s.DropQueue {
		adjustedStats, ok := s.ArmyPool[drop.TroopID]
		if !ok {
			continue
		}
		s.Troops = append(s.Troops, TroopState{
			ID:          s.NextTroopID,
			TroopID:     drop.TroopID,
			X:           float64(drop.X),
			Y:           float64(drop.Y),
			CurrentHP:   adjustedStats.Defence,
			MaxHP:       adjustedStats.Defence,
			Attack:      adjustedStats.TroopAttackPower,
			Range:       float64(adjustedStats.Range),
			Speed:       float64(adjustedStats.MovementSpeed) * 0.05,
			Name:        adjustedStats.Name,
			ElementType: adjustedStats.AttributeStrength,
		})
		s.NextTroopID++
	}
	s.DropQueue = nil

	// troops move toward nearest building or attack if in range
	for i := range s.Troops {
		target := nearestBuilding(s.Troops[i], s.Buildings)
		if target == nil {
			continue
		}

		center := buildingCenter(*target)
		troopPos := Point{X: s.Troops[i].X, Y: s.Troops[i].Y}
		dist := distance(troopPos, center)

		if dist <= s.Troops[i].Range {
			damage := calculateDamage(
				s.Troops[i].Attack,
				s.Troops[i].ElementType,
				target.ElementType,
			)
			target.CurrentHP -= damage
		} else {
			dx := center.X - s.Troops[i].X
			dy := center.Y - s.Troops[i].Y
			length := math.Sqrt(dx*dx + dy*dy)
			s.Troops[i].X += (dx / length) * s.Troops[i].Speed
			s.Troops[i].Y += (dy / length) * s.Troops[i].Speed
		}
	}

	// buildings AOE attack all troops in range
	for i := range s.Buildings {
		if s.Buildings[i].CurrentHP <= 0 {
			continue
		}
		targets := troopsInRange(s.Buildings[i], s.Troops)
		for _, t := range targets {
			damage := calculateDamage(
				s.Buildings[i].Attack,
				s.Buildings[i].ElementType,
				t.ElementType,
			)
			t.CurrentHP -= damage/40
		}
	}

	// remove dead troops
	aliveTroops := []TroopState{}
	for _, t := range s.Troops {
		if t.CurrentHP > 0 {
			aliveTroops = append(aliveTroops, t)
		}
	}
	s.Troops = aliveTroops

	// remove dead buildings
	aliveBuildings := []BuildingState{}
	for _, b := range s.Buildings {
		if b.CurrentHP > 0 {
			aliveBuildings = append(aliveBuildings, b)
		}
	}
	s.Buildings = aliveBuildings

	// check battle end
	if len(s.Buildings) == 0 || (len(s.Troops) == 0 && len(s.DropQueue) == 0 && s.ClientConnected && s.ElapsedTicks > 200) {
		s.Done = true
		log.Printf("battle stopped on purpose")
	}
}

func (s *BattleSession) GetState() BattleState {
	s.mu.Lock()
	defer s.mu.Unlock()

	destroyed := s.InitBuildings - len(s.Buildings)
	destruction := 0
	if s.InitBuildings > 0 {
		destruction = (destroyed * 100) / s.InitBuildings
	}

	return BattleState{
		Buildings:    s.Buildings,
		Troops:       s.Troops,
		Done:         s.Done,
		Destruction:  destruction,
		ElapsedTicks: s.ElapsedTicks,
	}
}

func (s *BattleSession) DropTroop(drop TroopDrop) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DropQueue = append(s.DropQueue, drop)
}