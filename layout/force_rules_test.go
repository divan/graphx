package layout

import (
	"testing"

	"github.com/divan/graphx/graph"
)

func TestEachOnEach(t *testing.T) {
	objects := make(map[string]*Object)

	objects["1"] = NewObjectID(NewPoint(1, 1, 1), "1")
	objects["2"] = NewObjectID(NewPoint(2, 2, 2), "2")
	objects["3"] = NewObjectID(NewPoint(3, 3, 3), "3")

	links := []*graph.Link{}
	force := NewGravityForce(-10, nil)

	EachOnEach(force, objects, links)
}
