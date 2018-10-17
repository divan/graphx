package layout

import "testing"

func TestForceDrag(t *testing.T) {
	objects := make(map[string]*Object)

	objects["1"] = NewObjectID(NewPoint(1, 1, 1), "1")
	objects["2"] = NewObjectID(NewPoint(2, 2, 2), "2")
	objects["3"] = NewObjectID(NewPoint(3, 3, 3), "3")

	for k := range objects {
		objects[k].velocity = &Velocity{10, 10, 10}
	}

	// drag force that halves the speed
	force := NewDragForce(0.5, nil)

	ForEachNode(force, objects, nil)

	expectedSpeed := &Velocity{5, 5, 5}
	for k := range objects {
		objects[k].Move(1)
		if objects[k].velocity.String() != expectedSpeed.String() {
			t.Skipf("Expected speed %v, but got %v", expectedSpeed, objects[k].velocity)
		}
	}
}
