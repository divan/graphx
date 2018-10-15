package layout

import (
	"testing"

	"github.com/divan/graphx/generation/basic"
	"github.com/divan/graphx/graph"
)

func TestLayout(t *testing.T) {
	graph := basic.NewLineGenerator(2).Generate()
	repelling := NewGravityForce(-1.0, BarneHutMethod)
	//springs := NewSpringForce(0.01, 12.0, ForEachLink)
	l := New(graph, repelling)
	l.objects["0"].SetPosition(0, 0, 0)
	l.objects["1"].SetPosition(100, 100, 100)

	for i := 0; i < 1000; i++ {
		l.UpdatePositions()
	}
}

func TestLayoutAdd(t *testing.T) {
	g := graph.NewGraph()
	g.AddNode(graph.NewBasicNode("node 0"))
	g.AddNode(graph.NewBasicNode("node 1"))
	g.AddLink("node 0", "node 1")
	repelling := NewGravityForce(-1.0, ForEachNode)
	l := New(g, repelling)

	if len(l.objects) != 2 {
		t.Fatalf("objects map expected to be of %d length, but is of %d", 2, len(l.objects))
	}
	if len(l.positions) != 2 {
		t.Fatalf("positions expected to be of %d length, but is of %d", 2, len(l.objects))
	}
	if len(l.links) != 1 {
		t.Fatalf("expected to have %d links, but has %d", 1, len(l.links))
	}

	l.AddNode(graph.NewBasicNode("node 2"))
	if len(l.objects) != 3 {
		t.Fatalf("objects map expected to be of %d length, but is of %d", 3, len(l.objects))
	}
	if len(l.positions) != 3 {
		t.Fatalf("positions expected to be of %d length, but is of %d", 3, len(l.objects))
	}
}
