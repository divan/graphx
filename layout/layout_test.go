package layout

import (
	"testing"

	"github.com/divan/graphx/generation/basic"
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
