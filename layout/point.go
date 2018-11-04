package layout

import (
	"fmt"
	"math"

	"github.com/divan/graphx/layout/octree"
)

// Point represents a single point in 3D space.
type Point struct {
	_X float64 `json:"x"`
	_Y float64 `json:"y"`
	_Z float64 `json:"z"`
}

// NewPoint creates a new points for the given coordinates.
func NewPoint(x, y, z float64) *Point {
	return &Point{
		_X: x,
		_Y: y,
		_Z: z,
	}
}

// SetPosition sets points positon to the given coordines.
func (p *Point) SetPosition(x, y, z float64) {
	p._X = x
	p._Y = y
	p._Z = z
}

// distance calculated distance betweein two objects in 3D space.
func distance(from, to octree.Point) float64 {
	dx := float64(to.X() - from.X())
	dy := float64(to.Y() - from.Y())
	dz := float64(to.Z() - from.Z())
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// String implements Stringer for Point.
func (p Point) String() string {
	return fmt.Sprintf("(%+3.3f,%+3.3f,%+3.3f)", p._X, p._Y, p._Z)
}

// Implement octree.Point
func (p Point) X() float64 { return p._X }
func (p Point) Y() float64 { return p._Y }
func (p Point) Z() float64 { return p._Z }
