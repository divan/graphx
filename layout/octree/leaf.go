package octree

import "fmt"

// Leaf represents Octant without children, "external node". Satisfies Octant.
type Leaf struct {
	Point

	octree *Octree // to access cache of octree
}

// make sure Leaf satisfies Octant interface at compile time.
var _ = Octant(&Leaf{})

// NewLeaf initializes a new Leaf.
func (o *Octree) NewLeaf(p Point) *Leaf {
	leaf := &Leaf{
		Point:  p,
		octree: o,
	}
	o.ids[p.ID()] = leaf
	return leaf
}

// Center returns point of the Leaf. Implements Octant interface.
func (l *Leaf) Center() Point {
	return l.Point
}

// Insert inserts new Point into existing Leaf and returns updated
// node, which may be transformed into node. Implements Octant interface.
func (l *Leaf) Insert(p Point) Octant {
	if l == nil {
		return l.octree.NewLeaf(p)
	}

	// remove cached leaf id
	delete(l.octree.ids, p.ID())

	//external node, and we have two points in one Octant.
	//need to convert it to internal node and divide
	n := l.octree.NewNode()
	n.massCenter = massCenter(l.Center(), p)
	n.Insert(l.Center())
	n.Insert(p)

	return n
}

// String implements fmt.Stringer interface for Leaf.
func (l *Leaf) String() string {
	if l == nil || l.Center() == nil {
		return "."
	}
	c := l.Center()
	return fmt.Sprintf("L: [%.1f, %.1f, %.1f]", c.X(), c.Y(), c.Z())
}
