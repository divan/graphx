package octree

// Point represents a point in 3D space with a mass.
type Point interface {
	ID() string
	X() float64
	Y() float64
	Z() float64
	Mass() float64
}

// point implements a Point for mass usage with centers.
type point struct {
	id      string
	x, y, z float64
	mass    float64
}

// NewPoint creates a new point for the given coords and mass.
func NewPoint(id string, x, y, z, m float64) Point {
	return &point{
		id:   id,
		x:    x,
		y:    y,
		z:    z,
		mass: m,
	}
}

func (p *point) ID() string    { return p.id }
func (p *point) X() float64    { return p.x }
func (p *point) Y() float64    { return p.y }
func (p *point) Z() float64    { return p.z }
func (p *point) Mass() float64 { return p.mass }
