package layout

import (
	"fmt"
	"math"
)

// Object represents an object in 3D space with some ID information
// attached to it.
type Object struct {
	ID string

	X, Y, Z int
	Mass    int

	velocity *Velocity
	force    *ForceVector
}

// NewObject creates new object with given coordinates.
func NewObject(x, y, z int) *Object {
	return &Object{
		X:    x,
		Y:    y,
		Z:    z,
		Mass: 1,

		velocity: ZeroVelocity,
		force:    ZeroForce,
	}
}

// NewObjectID creates new object with given coordinates and ID.
func NewObjectID(x, y, z int, id string) *Object {
	ret := NewObject(x, y, z)
	ret.ID = id
	return ret
}

// String implements Stringer interface for Object.
func (o *Object) String() string {
	return fmt.Sprintf("[%d, %d, %d, m: %d]", o.X, o.Y, o.Z, o.Mass)
}

// Move updates object positions by calculating movement with current force and
// velocity in a time interval dt.
func (o *Object) Move(dt int) (dx, dy, dz float64) {
	o.updateVelocity(dt)
	dx, dy, dz = o.velocity.Distance(dt)
	o.X += int(dx)
	o.Y += int(dy)
	o.Z += int(dz)
	return
}

// updateVelocity updates object velocity with a current force applied.
func (o *Object) updateVelocity(dt int) {
	if o.force == nil {
		return
	}

	o.velocity.X += float64(dt) * o.force.DX / float64(o.Mass)
	o.velocity.Y += float64(dt) * o.force.DY / float64(o.Mass)
	o.velocity.Z += float64(dt) * o.force.DZ / float64(o.Mass)
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
var ZeroVelocity = &Velocity{}

// Distance calculates distance for the given time interval dt for three dimentions.
func (v *Velocity) Distance(dt int) (dx, dy, dz float64) {
	dx = float64(dt) * v.X
	dy = float64(dt) * v.Y
	dz = float64(dt) * v.Z
	return
}
