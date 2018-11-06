package octree

import (
	"fmt"
)

// Node represents Octant with children, "internal node". Satisifies Octant.
type Node struct {
	Leafs      *[8]Octant
	massCenter Point
	box        *Box

	octree *Octree // to access cache of octree
}

// make sure Node satisfies Octant interface at compile time.
var _ = Octant(&Node{})

// NewNode initializes a new Node.
func (o *Octree) NewNode(box *Box) *Node {
	var leafs [8]Octant
	return &Node{
		Leafs:  &leafs,
		octree: o,
		box:    box,
	}
}

// Center returns center of the mass of the Node. Implements Octant interface.
func (n *Node) Center() Point {
	return n.massCenter
}

// Insert inserts new Point into existing Node and returns
// updated Node. Implements Octant interface.
func (n *Node) Insert(p Point, parentBox *Box) Octant {
	idx, box := n.findOctantIdx(p, *parentBox)
	leaf := n.Leafs[idx]
	var l Octant
	if leaf == nil {
		l = n.octree.NewLeaf(p, box)
	} else {
		l = leaf.Insert(p, box)
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
// given Octant and its Box. It's in following order:
// 0 - Top, Front, Left
// 1 - Top, Front, Right
// 2 - Top, Back, Left
// 3 - Top, Back, Right
// 4 - Bottom, Front, Left
// 5 - Bottom, Front, Right
// 6 - Bottom, Back, Left
// 7 - Bottom, Back, Right
func (n *Node) findOctantIdx(p Point, parentBox Box) (int, *Box) {
	midX, midY, midZ := parentBox.Center()

	var (
		i   int
		box = parentBox // build new box based on parentBox
	)
	if p.X() > midX {
		i |= 1

		box.Left = midX
	} else {
		box.Right = midX
	}

	if p.Y() > midY {
		i |= 2

		box.Bottom = midY
	} else {
		box.Top = midY
	}

	if p.Z() > midZ {
		i |= 4
		box.Front = midZ
	} else {
		box.Back = midZ
	}
	return i, &box
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
	return n.box.Width()
}

// Box returns bounding box. Implements Octant.
func (n *Node) Box() *Box {
	return n.box
}
