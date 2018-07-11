package basic

import "testing"

func TestGrid3D(t *testing.T) {
	rows, cols, levels := 3, 3, 3
	gen := NewGrid3DGenerator(rows, cols, levels)
	g := gen.Generate()

	got := len(g.Nodes())
	expected := rows * cols * levels
	if got != expected {
		t.Fatalf("Expected graph to have %d nodes, but got %d", expected, got)
	}

	got = len(g.Links())
	expected = 54
	if got != expected {
		t.Fatalf("Expected graph to have %d links, but got %d", expected, got)
	}
}
