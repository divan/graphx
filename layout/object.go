package layout

import (
	"fmt"
)

// Object represents an object in 3D space with some ID information
// attached to it.
type Object struct {
	ID string
	*Point
	Mass float64

	velocity *Velocity
	force    *ForceVector
}

// NewObject creates new object with given point.
func NewObject(point *Point) *Object {
	return &Object{
		Point: point,
		Mass:  1,

		velocity: ZeroVelocity(),
		force:    ZeroForce(),
	}
}

// NewObjectID creates new object with given coordinates and ID.
func NewObjectID(p *Point, id string) *Object {
	ret := NewObject(p)
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
