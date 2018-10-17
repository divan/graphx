package layout

import (
	"errors"
	"fmt"
	"math"
)

const theta = 0.5 // barne-hut defaults

// Octree represents Octree data structure.
// See https://en.wikipedia.org/wiki/Octree for details.
type Octree struct {
	root  octant
	force Force
}

// octant represent a node in octree, which is an octant of a cube.
// See: http://en.wikipedia.org/wiki/Octant_(solid_geometry)
type octant interface {
	Center() *Object
	Insert(p *Object) octant
}

// node represents octant with children, "internal node". Satisifies octant.
type node struct {
	leafs      *[8]octant
	massCenter *Object
}

// Center returns center of the mass of the node. Implements octant interface.
func (n *node) Center() *Object {
	return n.massCenter
}

// make sure node satisfies octant interface at compile time.
var _ = octant(&node{})

// leaf represents octant without children, "external node". Satisfies octant.
type leaf struct {
	point *Object
}

// Center returns point of the leaf. Implements octant interface.
func (l *leaf) Center() *Object {
	return l.point
}

// make sure leaf satisfies octant interface at compile time.
var _ = octant(&leaf{})

// NewOctree inits new octree.
func NewOctree(force Force) *Octree {
	return &Octree{
		force: force,
	}
}

// NewOctreeFromNodes inits new octree with current
// positions of the nodes and sets gravity force to force.
func NewOctreeFromNodes(objects map[string]*Object, force Force) *Octree {
	ot := NewOctree(force)
	for _, o := range objects {
		ot.Insert(o)
	}

	return ot
}

// newNode initializes a new node.
func newNode() *node {
	var leafs [8]octant
	for i := 0; i < 8; i++ {
		leafs[i] = newLeaf(nil)
	}
	return &node{
		leafs: &leafs,
	}
}

// newLeaf initializes a new leaf.
func newLeaf(p *Object) *leaf {
	return &leaf{
		point: p,
	}
}

// Insert adds new Point into the Octree data structure.
func (o *Octree) Insert(p *Object) {
	if o.root == nil {
		o.root = newLeaf(p)
		return
	}

	o.root = o.root.Insert(p)
}

// Insert inserts new Point into existing node and returns
// updated node. Implements octant interface.
func (n *node) Insert(p *Object) octant {
	idx := findOctantIdx(n, p)
	n.leafs[idx] = n.leafs[idx].Insert(p)
	n.updateMassCenter()
	return n
}

// Insert inserts new Point into existing leaf and returns updated
// node, which may be transformed into node. Implements octant interface.
func (l *leaf) Insert(p *Object) octant {
	if l == nil || l.Center() == nil {
		l = newLeaf(p)
		return l
	}

	//external node, and we have two points in one octant.
	//need to convert it to internal node and divide
	node := newNode()
	node.massCenter = massCenter([]*Object{l.Center(), p})

	idx1 := findOctantIdx(node, l.Center())
	idx2 := findOctantIdx(node, p)
	node.leafs[idx1] = l
	node.leafs[idx2] = newLeaf(p)
	return node
}

// update center of the mass of the given node, calculating it from
// leaf centers of the mass.
func (n *node) updateMassCenter() {
	var points []*Object
	for i := range n.leafs {
		points = append(points, n.leafs[i].Center())
	}

	n.massCenter = massCenter(points)
}

func massCenter(points []*Object) *Object {
	var (
		xm, ym, zm float64
		totalMass  float64
	)

	for _, p := range points {
		if p == nil {
			continue
		}
		totalMass += p.Mass
		xm += p.X * p.Mass
		ym += p.Y * p.Mass
		zm += p.Z * p.Mass
	}

	point := NewPoint(xm/totalMass, ym/totalMass, zm/totalMass)
	ret := NewObject(point)
	ret.Mass = totalMass
	return ret
}

// findOctantIdx returns index of 8-length array with children of the
// given octant. It's in following order:
// 0 - Top, Front, Left
// 1 - Top, Front, Right
// 2 - Top, Back, Left
// 3 - Top, Back, Right
// 4 - Bottom, Front, Left
// 5 - Bottom, Front, Right
// 6 - Bottom, Back, Left
// 7 - Bottom, Back, Right
func findOctantIdx(o octant, p *Object) int {
	center := o.Center()

	var i int
	if p.X > center.X {
		i |= 1
	}

	if p.Y > center.Y {
		i |= 2
	}

	if p.Z > center.Z {
		i |= 4
	}
	return i
}

// String implements Stringer interface for octree.
func (o *Octree) String() string {
	return fmt.Sprintf("Root: %T, leafs: %v", o.root, o.root.(*node).leafs)
}
func (n *node) String() string {
	var out string
	for i := 0; i < 8; i++ {
		if n.leafs[i] == nil {
			out += "."
		} else if l, ok := n.leafs[i].(*leaf); ok {
			if l == nil || l.Center() == nil {
				out += "."
			} else {
				out += "L"
			}
		} else if _, ok := n.leafs[i].(*node); ok {
			out += "N"
		}
	}
	return fmt.Sprintf("Node: (%.1f, %.1f, %.1f): [%s]", n.Center().X, n.Center().Y, n.Center().Z, out)
}

func (l *leaf) String() string {
	if l == nil || l.Center() == nil {
		return "."
	}
	c := l.Center()
	return fmt.Sprintf("L %s: [%.1f, %.1f, %.1f]", c.ID, c.X, c.Y, c.Z)
}

// CalcForce calculates force between two nodes using Barne-Hut method.
func (o *Octree) CalcForce(id string) (*ForceVector, error) {
	from, err := findLeaf(o.root, id)
	if err != nil {
		return nil, err
	}
	return o.calcForce(from, o.root), nil
}

func (o *Octree) calcForce(from *leaf, to octant) *ForceVector {
	if from == nil {
		panic(errors.New("calcForce from nil"))
	}
	ret := ZeroForce()
	if toLeaf, ok := to.(*leaf); ok {
		if toLeaf == nil || toLeaf.Center() == nil {
			return ret
		}
		return o.force.Apply(from.Center(), toLeaf.Center())
	} else if toNode, ok := to.(*node); ok {
		// calculate ratio
		width := toNode.width()

		r := distance(from.Center().Point, to.Center().Point)

		if float64(width)/float64(r) < theta {
			return o.force.Apply(from.Center(), to.Center())
		}

		for i := range toNode.leafs {
			f := o.calcForce(from, toNode.leafs[i])
			ret.Add(f)
		}
	}
	return ret
}

// findLeaf finds leaf for Point by given id.
func findLeaf(o octant, id string) (*leaf, error) {
	if l, ok := o.(*leaf); ok {
		if l == nil || l.Center() == nil {
			return nil, errors.New("nil leaf")
		}
		if l.point.ID == id {
			return l, nil
		}
	} else if n, ok := o.(*node); ok {
		for i := 0; i < 8; i++ {
			l1, err := findLeaf(n.leafs[i], id)
			if err == nil {
				return l1, nil
			}
		}
	}
	return nil, fmt.Errorf("node '%s' not found in octree", id)
}

// width returns width of the node, calculated from leaf coordinates.
func (n *node) width() int32 {
	// find two non-nil nodes
	for i := 0; i < 8; i++ {
		if n.leafs[i] != nil && n.leafs[i].Center() != nil {
			for j := 0; j < 8; j++ {
				if n.leafs[j] != nil && n.leafs[j].Center() != nil {
					p1, p2 := n.leafs[i].Center(), n.leafs[j].Center()
					// calculate non-zero difference in one of the dimensions (any)
					xwidth := math.Abs(float64(p1.X - p2.X))
					if xwidth > 0 {
						return int32(xwidth)
					}
					ywidth := math.Abs(float64(p1.Y - p2.Y))
					if ywidth > 0 {
						return int32(xwidth)
					}
					zwidth := math.Abs(float64(p1.Z - p2.Z))
					if zwidth > 0 {
						return int32(xwidth)
					}
				}
			}
		}
	}
	return 0
}
