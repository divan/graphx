package layout

import (
	"fmt"
	"math"
)

// Object represents an object in 3D space with some ID information
// attached to it.
type Object struct {
	ID string

	X, Y, Z float64
	Mass    float64

	velocity *Velocity
	force    *ForceVector
}

// NewObject creates new object with given coordinates.
func NewObject(x, y, z float64) *Object {
	return &Object{
		X:    x,
		Y:    y,
		Z:    z,
		Mass: 1,

		velocity: ZeroVelocity(),
		force:    ZeroForce(),
	}
}

// NewObjectID creates new object with given coordinates and ID.
func NewObjectID(x, y, z float64, id string) *Object {
	ret := NewObject(x, y, z)
	ret.ID = id
	return ret
}

// String implements Stringer interface for Object.
func (o *Object) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f, m: %.2f]", o.X, o.Y, o.Z, o.Mass)
}

// Move updates object positions by calculating movement with current force and
// velocity in a time interval dt.
func (o *Object) Move(dt int) (dx, dy, dz float64) {
	o.updateVelocity(dt, o.force)
	v := o.velocity
	t := float64(dt)
	o.X += t * v.X
	o.Y += t * v.Y
	o.Z += t * v.Z
	return o.X, o.Y, o.Z
}

// updateVelocity updates object velocity with a current force applied.
func (o *Object) updateVelocity(dt int, force *ForceVector) {
	if o.force == ZeroForce() {
		return
	}

	o.velocity.X += float64(dt) * force.DX / float64(o.Mass)
	o.velocity.Y += float64(dt) * force.DY / float64(o.Mass)
	o.velocity.Z += float64(dt) * force.DZ / float64(o.Mass)
}

// SetPosition sets object positon to the given coordines.
func (o *Object) SetPosition(x, y, z float64) {
	o.X = x
	o.Y = y
	o.Z = z
}

// distance calculated distance betweein two objects in 3D space.
func distance(from, to *Object) float64 {
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	dz := float64(to.Z - from.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

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
