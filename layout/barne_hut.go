package layout

import (
	"github.com/divan/graphx/layout/octree"
)

// DefaultTheta sets the default theta value (normally in 0.5-1.0 range)
const DefaultTheta = 0.5

// BarneHut implements n-body problem Barne-Hut optimization technicque.
type BarneHut struct {
	theta float64
	force Force // repelling force to use
}

// NewBarneHutMethod initializes a new Barne-Hut method helper.
func NewBarneHut(force Force) *BarneHut {
	return &BarneHut{
		theta: DefaultTheta,
		force: force,
	}
}

func (b *BarneHut) Calculate(objects map[string]*Object) map[string]*ForceVector {
	oc := octree.New()
	for _, o := range objects {
		oc.Insert(o)
	}

	forces := make(map[string]*ForceVector)
	for id := range objects {
		forces[id] = b.CalcForce(oc, id)
	}

	return forces
}

// CalcForce calculates force between two nodes using Barne-Hut method.
func (b *BarneHut) CalcForce(oc *octree.Octree, id string) *ForceVector {
	from, err := oc.FindLeaf(id)
	if err != nil {
		panic(err)
	}
	return b.calcForce(oc, from, oc.Root)
}

func (b *BarneHut) calcForce(oc *octree.Octree, from, to octree.Octant) *ForceVector {
	ret := ZeroForce()
	if leaf, ok := to.(*octree.Leaf); ok {
		if leaf == nil {
			return ret
		}
		return b.force.Apply(from.Center().(*Object), to.Center().(*Object))
	} else if node, ok := to.(*octree.Node); ok {
		// calculate ratio
		width := node.Width()

		r := distance(from.Center(), to.Center())

		if width/r < b.theta {
			// FIXME. TODO(divan): this is temporary hack to verify tests.
			// Refactor point/object represtnations here and in octree and/or refactor forces
			c := from.Center()
			f, ok := c.(*Object)
			if !ok {
				f = NewObject(c.X(), c.Y(), c.Z())
			}
			c = to.Center()
			t, ok := c.(*Object)
			if !ok {
				t = NewObject(c.X(), c.Y(), c.Z())
			}
			return b.force.Apply(f, t)
		}

		for _, l := range node.Leafs {
			if l == nil {
				continue
			}
			f := b.calcForce(oc, from, leaf)
			ret.Add(f)
		}
	}
	return ret
}
