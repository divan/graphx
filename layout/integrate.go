package layout

import (
	"math"
)

// integrate performs forces integration using Euler numerical
// integration method.
//
// F = d(m * v) / dt
//  (mass is constant in our case)
// v = d{x,y,z}/dt
//
// dv = dt * F / m
//
// d{x,y,z} = v * dt
//
// returns total movement amount
func (l *Layout) integrate() float64 {
	const dt = 3           // integration step
	var tx, ty, tz float64 // total movement
	for k := range l.objects {

		body := l.objects[k]
		dvx, dvy, dvz := body.velocity.X, body.velocity.Y, body.velocity.Z
		v := math.Sqrt(dvx*dvx + dvy*dvy + dvz*dvz)

		if v > 1 {
			body.velocity.X = dvx / v
			body.velocity.Y = dvy / v
			body.velocity.Z = dvz / v
		}

		dx, dy, dz := body.Move(dt)

		// calculate total displacement
		tx += math.Abs(dx)
		ty += math.Abs(dy)
		tz += math.Abs(dz)
	}

	return (tx*tx + ty*ty + tz*tz) / float64(len(l.objects))
}
