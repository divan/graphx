package octree

import "testing"

func TestNodeWidth(t *testing.T) {
	p1 := NewPoint("0", 100, 0, 0, 1)
	p2 := NewPoint("1", -100, 0, 0, 1)
	octree := New(NewBoxPoints(p1, p2))
	octree.Insert(p1)
	octree.Insert(p2)
	root := octree.Root.(*Node)
	expected := 200.0
	got := root.Width()
	if got != expected {
		t.Fatalf("Expect width to be %f, but got %f", expected, got)
	}

	octree.Insert(NewPoint("2", 80, 0, 0, 1))
	l := root.Leafs[1].(*Node)
	expected = 100.0
	got = l.Width()
	if got != expected {
		t.Fatalf("Expect width to be %f, but got %f", expected, got)
	}
}

func TestFindOctant(t *testing.T) {
	p1 := NewPoint("0", 1, 1, 1, 1)
	p2 := NewPoint("1", 100, 100, 100, 1)
	octree := New(NewBoxPoints(p1, p2))
	octree.Insert(p1)
	octree.Insert(p2)
	node := octree.Root.(*Node)
	box := *(node.Box())
	p := NewPoint("", 90, 90, 90, 1)
	idx, _ := node.findOctantIdx(p, box)
	expected := 7 // Bottom, Back, Right
	if idx != expected {
		t.Fatalf("Expect idx to be %d, but got %d", expected, idx)
	}
	p = NewPoint("", 10, 10, 10, 1)
	idx, _ = node.findOctantIdx(p, box)
	expected = 0 // Top, Front, Left
	if idx != expected {
		t.Fatalf("Expect idx to be %d, but got %d", expected, idx)
	}
}
