package octree

import (
	"errors"
	"fmt"
)

// Octree represents Octree data structure.
// See https://en.wikipedia.org/wiki/Octree for details.
type Octree struct {
	Root Octant
}

// octant represent a node in octree, which is an octant of a cube.
// See: http://en.wikipedia.org/wiki/Octant_(solid_geometry)
type Octant interface {
	Center() Point
	Insert(p Point) Octant
}

// New inits new octree.
func New() *Octree {
	return &Octree{}
}

// Insert adds new Point into the Octree data structure.
func (o *Octree) Insert(p Point) {
	if o.Root == nil {
		o.Root = NewLeaf(p)
		return
	}

	o.Root = o.Root.Insert(p)
}

// FindLeafs searches for the leaf with the given id.
func (o *Octree) FindLeaf(id string) (Octant, error) {
	oct, ok := o.findLeaf(o.Root, id)
	if !ok {
		return nil, errors.New("leaf not found")
	}
	return oct, nil
}

func (o *Octree) findLeaf(oct Octant, id string) (Octant, bool) {
	switch x := oct.(type) {
	case *Leaf:
		if x.ID() == id {
			return x, true
		}
		return nil, false
	case *Node:
		for i := 0; i < 8; i++ {
			if x.Leafs[i] == nil {
				continue
			}
			leaf, ok := o.findLeaf(x.Leafs[i], id)
			if ok {
				return leaf, true
			}
		}
	}
	return nil, false
}

// String implements Stringer interface for octree.
func (o *Octree) String() string {
	return fmt.Sprintf("Root: %T, leafs: %v", o.Root, o.Root.(*Node).Leafs)
}
