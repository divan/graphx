package layout

import (
	"fmt"
	"math"

	"github.com/divan/graphx/layout/octree"
)

// Object represents an object in 3D space with some ID information
// attached to it.
type Object struct {
	_ID   string  `json:"-"`
	_X    float64 `json:"x"`
	_Y    float64 `json:"y"`
	_Z    float64 `json:"z"`
	_Mass float64 `json:"mass"`

	velocity *Velocity    `json:"-"`
	force    *ForceVector `json:"-"`
}

// NewObject creates new object with given point.
func NewObject(x, y, z float64) *Object {
	return &Object{
		_X:    x,
		_Y:    y,
		_Z:    z,
		_Mass: 1,

		velocity: ZeroVelocity(),
		force:    ZeroForce(),
	}
}

// NewObjectID creates new object with given coordinates and ID.
func NewObjectID(x, y, z float64, id string) *Object {
	ret := NewObject(x, y, z)
	ret._ID = id
	return ret
}

// String implements Stringer interface for Object.
func (o *Object) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f, m: %.2f]", o.X(), o.Y(), o.Z(), o.Mass())
}

// X implements Point interface.
func (o *Object) X() float64 { return o._X }

// Y implements Point interface.
func (o *Object) Y() float64 { return o._Y }

// Z implements Point interface.
func (o *Object) Z() float64 { return o._Z }

// ID implements Point interface.
func (o *Object) ID() string { return o._ID }

// Mass implements Point interface.
func (o *Object) Mass() float64 { return o._Mass }

// SetPosition sets points positon to the given coordines.
func (o *Object) SetPosition(x, y, z float64) {
	o._X = x
	o._Y = y
	o._Z = z
}

// Move updates object positions by calculating movement with current force and
// velocity in a time interval dt.
func (o *Object) Move(dt int) (dx, dy, dz float64) {
	o.updateVelocity(dt, o.force)
	v := o.velocity
	t := float64(dt)
	o._X += t * v.X
	o._Y += t * v.Y
	o._Z += t * v.Z
	return o._X, o._Y, o._Z
}

// updateVelocity updates object velocity with a current force applied.
func (o *Object) updateVelocity(dt int, force *ForceVector) {
	if o.force == ZeroForce() {
		return
	}

	o.velocity.X += float64(dt) * force.DX / float64(o.Mass())
	o.velocity.Y += float64(dt) * force.DY / float64(o.Mass())
	o.velocity.Z += float64(dt) * force.DZ / float64(o.Mass())
}

func (o Object) Force() *ForceVector {
	return o.force
}
func (o Object) Velocity() *Velocity {
	return o.velocity
}

// distance calculated distance betweein two objects in 3D space.
func distance(from, to octree.Point) float64 {
	dx := float64(to.X() - from.X())
	dy := float64(to.Y() - from.Y())
	dz := float64(to.Z() - from.Z())
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
