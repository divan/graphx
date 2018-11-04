package layout

import "fmt"

// Velocity represents velocity vector.
type Velocity struct {
	X float64
	Y float64
	Z float64
}

// ZeroVelocity is, well, zero value for velocity.
func ZeroVelocity() *Velocity {
	return &Velocity{}
}

// String implements Stringer interface for velocity.
func (v *Velocity) String() string {
	return fmt.Sprintf("V(%.1f, %.1f, %.1f)", v.X, v.Y, v.Z)
}
