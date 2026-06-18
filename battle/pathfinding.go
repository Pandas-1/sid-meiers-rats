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

func inRange(b BuildingState, troop TroopState) bool {
    center := buildingCenter(b)
    troopPos := Point{X: troop.X, Y: troop.Y}
    return distance(troopPos, center) <= troop.Range
}