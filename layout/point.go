package layout

import "math"

// Point represents a single point in 3D space.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// NewPoint creates a new points for the given coordinates.
func NewPoint(x, y, z float64) *Point {
	return &Point{
		X: x,
		Y: y,
		Z: z,
	}
}

// SetPosition sets points positon to the given coordines.
func (p *Point) SetPosition(x, y, z float64) {
	p.X = x
	p.Y = y
	p.Z = z
}

// distance calculated distance betweein two objects in 3D space.
func distance(from, to *Point) float64 {
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	dz := float64(to.Z - from.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
