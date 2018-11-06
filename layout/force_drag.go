package layout

// DragForce represents general drag force with linear coefficient.
type DragForce struct {
	Coeff float64
	rule  ForceRule
}

// NewDragForce creates new drag force with the given drag coefficient.
func NewDragForce(coeff float64, rule ForceRule) Force {
	return &DragForce{
		Coeff: coeff,
		rule:  rule,
	}
}

// Apply calculates the drag force for `from` node.
// Second parameter is ignored. Satisfies Force interface.
// TODO(divan): find how to generalize force better, as here
// we don't need two nodes.
func (g *DragForce) Apply(p, _ Point) *ForceVector {
	pv, ok := p.(HasVelocity)
	if !ok {
		return ZeroForce()
	}
	return &ForceVector{
		DX: -g.Coeff * pv.Velocity().X,
		DY: -g.Coeff * pv.Velocity().Y,
		DZ: -g.Coeff * pv.Velocity().Z,
	}
}

// Name returns name of the force. Satisifies Force interface.
func (g *DragForce) Name() string {
	return "drag"
}

// Rule returns rule function to apply rules.
func (g *DragForce) Rule() ForceRule {
	return g.rule
}
