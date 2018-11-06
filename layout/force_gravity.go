package layout

import (
	"math"
)

// GravityForce represents gravity force,
// calculated by Coloumb's law. Implements Force interface.
type GravityForce struct {
	Coeff float64
	rule  ForceRule
}

// NewGravityForce creates new gravity force with the given
// Coloumb's law coefficient value.
func NewGravityForce(coeff float64, rule ForceRule) Force {
	return &GravityForce{
		Coeff: coeff,
		rule:  rule,
	}
}

// Apply calculates the gravity force between two nodes. Satisfies Force interface.
func (g *GravityForce) Apply(from, to Point) *ForceVector {
	xx := to.X() - from.X()
	yy := to.Y() - from.Y()
	zz := to.Z() - from.Z()

	r := math.Sqrt(xx*xx + yy*yy + zz*zz)
	if r == 0 {
		r = 10
	}

	v := g.Coeff * from.Mass() * to.Mass() / (r * r * r)
	return &ForceVector{
		DX: (xx * v),
		DY: (yy * v),
		DZ: (zz * v),
	}
}

// Name returns name of the force. Satisifies Force interface.
func (g *GravityForce) Name() string {
	return "gravity"
}

// Rule returns rule function to apply rules.
func (g *GravityForce) Rule() ForceRule {
	return g.rule
}
