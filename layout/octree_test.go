package layout

import "testing"

func TestInsert(t *testing.T) {
	o := NewOctree(NewGravityForce(-10, BarneHutMethod))
	p1 := NewObject(NewPoint(1, 1, 1))
	p1.Mass = 10
	o.Insert(p1)

	if o.root == nil {
		t.Fatalf("Expected root node to be non-nil")
	}

	center := o.root.Center()
	if center != p1 {
		t.Fatalf("Expected center to be %v, but got %v", p1, center)
	}

	p2 := NewObject(NewPoint(9, 9, 9))
	p2.Mass = 10
	o.Insert(p2)

	center = o.root.Center()
	expected := NewObject(NewPoint(5, 5, 5))
	expected.Mass = 20
	if center.String() != expected.String() {
		t.Fatalf("Expected center to be %v, but got %v", expected, center)
	}
}

func SkipFindOctantIdx(t *testing.T) {
	var tests = []struct {
		name string
		p    *Object
		idx  int
	}{
		{
			name: "bottom back right",
			p:    NewObject(NewPoint(9, 9, 9)),
			idx:  7,
		},
		{
			name: "top front left",
			p:    NewObject(NewPoint(1, 1, 1)),
			idx:  0,
		},
		{
			name: "bottom front right",
			p:    NewObject(NewPoint(9, 2, 9)),
			idx:  5,
		},
	}

	n := newNode()
	n.Insert(NewObject(NewPoint(5, 5, 5)))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			idx := n.findOctantIdx(test.p)
			if idx != test.idx {
				t.Fatalf("Expected idx %d, but got %d", test.idx, idx)
			}
		})
	}
}

func TestLeafInsert(t *testing.T) {
	p1 := NewObject(NewPoint(1, 1, 1))
	p2 := NewObject(NewPoint(-1, -1, -1))
	l := newLeaf(p1)
	center := l.Center()
	if center != p1 {
		t.Fatalf("center != p1")
	}
	node := l.Insert(p2)
	center = node.Center()
	expected := NewObject(NewPoint(0, 0, 0)) // zero coords
	expected.Mass = 2
	if center.String() != expected.String() {
		t.Fatalf("Expected %v, but got %v", expected, center)
	}
}

func TestBugCase1(t *testing.T) {
	o := NewOctree(NewGravityForce(-10, BarneHutMethod))
	objects := []*Object{
		NewObjectID(NewPoint(-2, 4, 1), "1"),
		NewObjectID(NewPoint(-6, 4, -1), "2"),
		NewObjectID(NewPoint(-1, -13, 3), "3"),
		NewObjectID(NewPoint(14, 14, 5), "4"),
		NewObjectID(NewPoint(-19, -5, 9), "5"),
	}
	for i := 0; i < 5; i++ {
		objects[i].Mass = 2
		o.Insert(objects[i])
	}
	for i := 0; i < 5; i++ {
		leaf, ok := o.root.FindLeaf(objects[i].ID)
		if !ok {
			t.Fatalf("Failed to find node %s", objects[i].ID)
		}
		if leaf.point.ID != objects[i].ID {
			t.Fatalf("Expected point index to be %s, got %s", objects[i].ID, leaf.point.ID)
		}
		if leaf.point != objects[i] {
			t.Fatalf("Expected point to be %v, got %v", objects[i], leaf.point)
		}
	}
}
