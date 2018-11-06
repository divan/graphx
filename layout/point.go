package layout

// Point represents a point in 3D space with a mass.
type Point interface {
	X() float64
	Y() float64
	Z() float64
	Mass() float64
}

// HasVelocity represents any point with velocity information.
type HasVelocity interface {
	Velocity() *Velocity
}
