package layout

import (
	"math"
	"testing"

	"github.com/divan/graphx/generation/basic"
)

func TestSpring(t *testing.T) {
	t.Run("2_bodies", func(t *testing.T) {
		restLength := 17.32
		graph := basic.NewLineGenerator(2).Generate()
		spring := NewSpringForce(0.01, restLength-2, ForEachLink) // -2 = let it oscillate a bit
		l := New(graph, spring)

		l.objects["0"].SetPosition(0, 0, 0)
		l.objects["1"].SetPosition(10, 10, 10)

		for i := 0; i < 1000; i++ {
			l.UpdatePositions()
			left, right := l.objects["0"], l.objects["1"]
			checkDistance(t, left, right, 0, 10)
		}
		left, right := l.objects["0"], l.objects["1"]
		d := distance(left.Point, right.Point)
		if math.Abs(d-restLength) > 5.0 {
			t.Logf("Expect diff to be less than %v, got %v", 5.0, math.Abs(d-restLength))
		}
	})
	t.Run("2_bodies_larger", func(t *testing.T) {
		restLength := 10.0
		graph := basic.NewLineGenerator(2).Generate()
		spring := NewSpringForce(0.02, restLength, ForEachLink)
		l := New(graph, spring)

		l.objects["0"].SetPosition(0, 0, 0)
		l.objects["1"].SetPosition(11, 0, 0)

		for i := 0; i < 100; i++ {
			l.UpdatePositions()
			left, right := l.objects["0"], l.objects["1"]
			checkDistanceX(t, left, right, 0, 10)
		}
		left, right := l.objects["0"], l.objects["1"]
		d := distance(left.Point, right.Point)
		if math.Abs(d-restLength) > 5.0 {
			t.Logf("Expect diff to be less than %v, got %v", 5.0, math.Abs(d-restLength))
		}
	})
}

// checkDistanceX checks left and right distances from their initial positions (x0 and x1) along X axis. they should be equal.
func checkDistanceX(t *testing.T, left, right *Object, x0, x1 float64) {
	dx1, dx2 := x0-left.X, right.X-x1
	if dx1-dx2 > 0.0001 {
		t.Fatalf("Expect dX be equal for left and right, but got (%.2f, %.2f)", dx1, dx2)
	}
}
