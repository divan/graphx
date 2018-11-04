package layout

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	var tests = []struct {
		from, to *Object
		expected float64
	}{
		{
			from:     NewObject(NewPoint(0, 0, 0)),
			to:       NewObject(NewPoint(10, 10, 10)),
			expected: math.Sqrt(300),
		},
		{
			from:     NewObject(NewPoint(2, 3, 1)),
			to:       NewObject(NewPoint(8, -5, 0)),
			expected: math.Sqrt(101),
		},
	}
	for _, test := range tests {
		got := distance(test.from, test.to)
		if got != test.expected {
			t.Fatalf("Expected %.3f, but got %.3f", test.expected, got)
		}
	}
}

func TestObjectVelocity(t *testing.T) {
	o := NewObject(NewPoint(1, 1, 1))
	o.velocity = &Velocity{10, 10, 10}

	o.updateVelocity(1, ZeroForce())
	expected := &Velocity{10, 10, 10}
	if o.velocity.String() != expected.String() {
		t.Fatalf("Expected %v, but got %v", expected, o.velocity)
	}

	o.updateVelocity(1, &ForceVector{-5, -5, -5})
	expected = &Velocity{5, 5, 5}
	if o.velocity.String() != expected.String() {
		t.Fatalf("Expected %v, but got %v", expected, o.velocity)
	}

	o.updateVelocity(1, &ForceVector{-5, -5, -5})
	expected = &Velocity{0, 0, 0}
	if o.velocity.String() != expected.String() {
		t.Fatalf("Expected %v, but got %v", expected, o.velocity)
	}

	o.velocity = &Velocity{10, 10, 10}
	o.updateVelocity(2, &ForceVector{-5, -5, -5})
	expected = &Velocity{0, 0, 0}
	if o.velocity.String() != expected.String() {
		t.Fatalf("Expected %v, but got %v", expected, o.velocity)
	}
}

func TestObjectMove(t *testing.T) {
	o := NewObject(NewPoint(100, 100, 100))
	o.velocity = ZeroVelocity() // &Velocity{10, 10, 10}
	o.force = &ForceVector{1, 1, 1}

	o.Move(3)
	t.Log(o.velocity)
	t.Log(o)
}
