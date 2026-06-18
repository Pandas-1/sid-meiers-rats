package battle

import "math"

type Point struct {
    X float64
    Y float64
}

func distance(p1, p2 Point) float64 {
    dx := p1.X - p2.X
    dy := p1.Y - p2.Y
    return math.Sqrt(dx*dx + dy*dy)
}

func buildingCenter(b BuildingState) Point {
    return Point{
        X: float64(b.X) + float64(b.Width)/2,
        Y: float64(b.Y) + float64(b.Height)/2,
    }
}

var elementChart = map[int]int{
    2: 5,  // Burning beats Ground
    3: 2,  // Wet beats Burning
    4: 3,  // Flying beats Wet
    5: 4,  // Ground beats Flying
    6: 7,  // Bright beats Dark
    7: 6,  // Dark beats Bright
}

func inRange(b BuildingState, troop TroopState) bool {
    center := buildingCenter(b)
    troopPos := Point{X: troop.X, Y: troop.Y}
    return distance(troopPos, center) <= troop.Range
}

func nearestBuilding(troop TroopState, buildings []BuildingState) *BuildingState {
    var nearest *BuildingState
    minDist := math.MaxFloat64

    for i := range buildings {
        center := buildingCenter(buildings[i])
        troopPos := Point{X: troop.X, Y: troop.Y}
        d := distance(troopPos, center)
        if d < minDist {
            minDist = d
            nearest = &buildings[i]
        }
    }
    return nearest
}

func troopsInRange(b BuildingState, troops []TroopState) []*TroopState {
    var inRangeTroops []*TroopState
    for i := range troops {
        if inRange(b, troops[i]) {
            inRangeTroops = append(inRangeTroops, &troops[i])
        }
    }
    return inRangeTroops
}

func calculateDamage(attack int, attackerElement int, defenderElement int) int {
    if attackerElement == 0 || defenderElement == 0 {
        return attack
    }

    strongAgainst, exists := elementChart[attackerElement]
    if !exists {
        return attack
    }

    if defenderElement == strongAgainst {
        // attacker strong against defender
        return int(float64(attack) * 1.375)
    }

    if elementChart[defenderElement] == attackerElement {
        // attacker weak against defender
        return int(float64(attack) * 0.5)
    }

    return attack
}
