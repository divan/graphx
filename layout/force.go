package layout

import "fmt"

// ForceVector represents the force vector in 3D space.
type ForceVector struct {
	DX float64 `json:"dx"`
	DY float64 `json:"dy"`
	DZ float64 `json:"dz"`
}

// Force defines the methods for physical force.
type Force interface {
	Name() string
	Apply(from, to *Object) *ForceVector
	Rule() ForceRule
}

// ZeroForce is a zero force.
func ZeroForce() *ForceVector { return &ForceVector{} }

// String implements Stringer interface for ForceVector.
func (f ForceVector) String() string {
	return fmt.Sprintf("f(%.03f, %.03f, %.03f)", f.DX, f.DY, f.DZ)
}

// Add adds new force to f.
func (f *ForceVector) Add(f1 *ForceVector) *ForceVector {
	f.DX += f1.DX
	f.DY += f1.DY
	f.DZ += f1.DZ
	return f
}

// Sub substracts new force from f.
func (f *ForceVector) Sub(f1 *ForceVector) *ForceVector {
	f.DX -= f1.DX
	f.DY -= f1.DY
	f.DZ -= f1.DZ
	return f
}

// Negative returns the same force vector, but in opposite direction.
func (f *ForceVector) Negative() *ForceVector {
	return &ForceVector{
		DX: -f.DX,
		DY: -f.DY,
		DZ: -f.DZ,
	}
}
