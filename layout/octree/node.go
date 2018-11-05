package octree

import (
	"fmt"
	"math"
)

// Node represents Octant with children, "internal node". Satisifies Octant.
type Node struct {
	Leafs      *[8]Octant
	massCenter Point
	width      float64 // box width

	octree *Octree // to access cache of octree
}

// make sure Node satisfies Octant interface at compile time.
var _ = Octant(&Node{})

// NewNode initializes a new Node.
func (o *Octree) NewNode() *Node {
	var leafs [8]Octant
	return &Node{
		Leafs:  &leafs,
		octree: o,
	}
}

// Center returns center of the mass of the Node. Implements Octant interface.
func (n *Node) Center() Point {
	return n.massCenter
}

// Insert inserts new Point into existing Node and returns
// updated Node. Implements Octant interface.
func (n *Node) Insert(p Point) Octant {
	idx := n.findOctantIdx(p)
	leaf := n.Leafs[idx]
	var l Octant
	if leaf == nil {
		l = n.octree.NewLeaf(p)
	} else {
		l = leaf.Insert(p)
	}

	n.Leafs[idx] = l
	n.updateWidth()
	return n
}

// update center of the mass of the given node, calculating it from
// leaf centers of the mass.
func (n *Node) updateMassCenter() {
	var points []Point
	for i := range n.Leafs {
		points = append(points, n.Leafs[i].Center())
	}

	n.massCenter = massCenter(points...)
}

func massCenter(points ...Point) Point {
	var (
		x, y, z float64
		mass    float64
	)

	for _, p := range points {
		if p == nil {
			continue
		}
		mass += p.Mass()
		x += p.X() * p.Mass()
		y += p.Y() * p.Mass()
		z += p.Z() * p.Mass()
	}

	return NewPoint("", x/mass, y/mass, z/mass, mass)
}

// findOctantIdx returns index of 8-length array with children of the
// given Octant. It's in following order:
// 0 - Top, Front, Left
// 1 - Top, Front, Right
// 2 - Top, Back, Left
// 3 - Top, Back, Right
// 4 - Bottom, Front, Left
// 5 - Bottom, Front, Right
// 6 - Bottom, Back, Left
// 7 - Bottom, Back, Right
func (n *Node) findOctantIdx(p Point) int {
	center := n.Center()

	var i int
	if p.X() > center.X() {
		i |= 1
	}

	if p.Y() > center.Y() {
		i |= 2
	}

	if p.Z() > center.Z() {
		i |= 4
	}
	return i
}

// String implements fmt.Stringer interface for Node.
func (n *Node) String() string {
	var out string
	for i := 0; i < 8; i++ {
		if n.Leafs[i] == nil {
			out += "."
		} else if l, ok := n.Leafs[i].(*Leaf); ok {
			if l == nil || l.Center() == nil {
				out += "."
			} else {
				out += "L"
			}
		} else if _, ok := n.Leafs[i].(*Node); ok {
			out += "N"
		}
	}
	return fmt.Sprintf("Node: (%.1f, %.1f, %.1f): [%s]", n.Center().X(), n.Center().Y(), n.Center().Z(), out)
}

// Width returns width of the Node, calculated from leaf coordinates.
func (n *Node) Width() float64 {
	return n.width
}

// updateWidth recalculates node's box width (useful for barne-hut method).
func (n *Node) updateWidth() {
	// find two non-nil nodes
	var max float64
	for i := 0; i < 8; i++ {
		if n.Leafs[i] != nil {
			for j := i + 1; j < 8; j++ {
				if n.Leafs[j] != nil {
					p1, p2 := n.Leafs[i].Center(), n.Leafs[j].Center()

					// calculate non-zero difference in one of the dimensions (any)
					dx := math.Abs(p1.X() - p2.X())
					dy := math.Abs(p1.Y() - p2.Y())
					dz := math.Abs(p1.Z() - p2.Z())

					width := math.Max(dx, math.Max(dy, dz))
					if width > max {
						max = width
					}
				}
			}
		}
	}
	n.width = max
}
