package octree

import (
	"fmt"
	"math"
)

// Node represents Octant with children, "internal node". Satisifies Octant.
type Node struct {
	Leafs      *[8]Octant
	massCenter Point
}

// make sure Node satisfies Octant interface at compile time.
var _ = Octant(&Node{})

// NewNode initializes a new Node.
func NewNode() *Node {
	var leafs [8]Octant
	return &Node{
		Leafs: &leafs,
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
		l = NewLeaf(p)
	} else {
		l = leaf.Insert(p)
	}

	n.Leafs[idx] = l
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
func (n *Node) Width() int32 {
	// find two non-nil nodes
	for i := 0; i < 8; i++ {
		if n.Leafs[i] != nil && n.Leafs[i].Center() != nil {
			for j := 0; j < 8; j++ {
				if n.Leafs[j] != nil && n.Leafs[j].Center() != nil {
					p1, p2 := n.Leafs[i].Center(), n.Leafs[j].Center()
					// calculate non-zero difference in one of the dimensions (any)
					xwidth := math.Abs(float64(p1.X() - p2.X()))
					if xwidth > 0 {
						return int32(xwidth)
					}
					ywidth := math.Abs(float64(p1.Y() - p2.Y()))
					if ywidth > 0 {
						return int32(xwidth)
					}
					zwidth := math.Abs(float64(p1.Z() - p2.Z()))
					if zwidth > 0 {
						return int32(xwidth)
					}
				}
			}
		}
	}
	return 0
}
