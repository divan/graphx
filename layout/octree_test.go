package layout

import "testing"

func TestInsert(t *testing.T) {
	o := NewOctree(NewGravityForce(-10, BarneHutMethod))
	p1 := NewObject(1, 1, 1)
	p1.Mass = 10
	o.Insert(p1)

	if o.root == nil {
		t.Fatalf("Expected root node to be non-nil")
	}

	center := o.root.Center()
	if center != p1 {
		t.Fatalf("Expected center to be %v, but got %v", p1, center)
	}

	p2 := NewObject(9, 9, 9)
	p2.Mass = 10
	o.Insert(p2)

	center = o.root.Center()
	expected := NewObject(5, 5, 5)
	expected.Mass = 20
	if center.String() != expected.String() {
		t.Fatalf("Expected center to be %v, but got %v", expected, center)
	}
}

func TestFindOctantIdx(t *testing.T) {
	var tests = []struct {
		name string
		p    *Object
		idx  int
	}{
		{
			name: "bottom back right",
			p:    NewObject(9, 9, 9),
			idx:  7,
		},
		{
			name: "top front left",
			p:    NewObject(1, 1, 1),
			idx:  0,
		},
		{
			name: "bottom front right",
			p:    NewObject(9, 2, 9),
			idx:  5,
		},
	}

	o := newLeaf(NewObject(5, 5, 5))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			idx := findOctantIdx(o, test.p)
			if idx != test.idx {
				t.Fatalf("Expected idx %d, but got %d", test.idx, idx)
			}
		})
	}
}

func TestLeafInsert(t *testing.T) {
	p1 := NewObject(1, 1, 1)
	p2 := NewObject(-1, -1, -1)
	l := newLeaf(p1)
	center := l.Center()
	if center != p1 {
		t.Fatalf("center != p1")
	}
	node := l.Insert(p2)
	center = node.Center()
	expected := NewObject(0, 0, 0) // zero coords
	expected.Mass = 2
	if center.String() != expected.String() {
		t.Fatalf("Expected %v, but got %v", expected, center)
	}
}

func TestBugCase1(t *testing.T) {
	o := NewOctree(NewGravityForce(-10, BarneHutMethod))
	objects := []*Object{
		NewObjectID(-2, 4, 1, "1"),
		NewObjectID(-6, 4, -1, "2"),
		NewObjectID(-1, -13, 3, "3"),
		NewObjectID(14, 14, 5, "4"),
		NewObjectID(-19, -5, 9, "5"),
	}
	for i := 0; i < 5; i++ {
		objects[i].Mass = 2
		o.Insert(objects[i])
	}
	for i := 0; i < 5; i++ {
		leaf, err := findLeaf(o.root, objects[i].ID)
		if err != nil {
			t.Fatalf("Expected err to be non nil, got %v", err)
		}
		if leaf.point.ID != objects[i].ID {
			t.Fatalf("Expected point index to be %s, got %s", objects[i].ID, leaf.point.ID)
		}
		if leaf.point != objects[i] {
			t.Fatalf("Expected point to be %v, got %v", objects[i], leaf.point)
		}
	}
}
