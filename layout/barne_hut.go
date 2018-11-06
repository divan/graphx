package layout

import (
	"math"

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
	box := boundingBox(objects)
	oc := octree.New(box)
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
	if to == from {
		return ret
	}
	if leaf, ok := to.(*octree.Leaf); ok {
		if leaf == nil {
			return ret
		}
		return b.force.Apply(from.Center(), to.Center())
	} else if node, ok := to.(*octree.Node); ok {
		// calculate ratio
		width := node.Width()

		r := distance(from.Center(), to.Center())

		if width/r < b.theta {
			return b.force.Apply(from.Center(), to.Center())
		}

		for _, l := range node.Leafs {
			if l == nil {
				continue
			}
			if from == l {
				continue
			}
			f := b.calcForce(oc, from, l)
			ret.Add(f)
		}
	}
	return ret
}

// boundingBox calculates the bounding box that fits all the points.
func boundingBox(objects map[string]*Object) *octree.Box {
	var (
		minX, maxX float64
		minY, maxY float64
		minZ, maxZ float64
	)
	for i := range objects {
		minX = math.Min(minX, objects[i].X())
		maxX = math.Max(maxX, objects[i].X())
		minY = math.Min(minY, objects[i].Y())
		maxY = math.Max(maxY, objects[i].Y())
		minZ = math.Min(minZ, objects[i].Z())
		maxZ = math.Max(maxZ, objects[i].Z())
	}
	return octree.NewBox(minX, maxX, minY, maxY, minZ, maxZ)
}
