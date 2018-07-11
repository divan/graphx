package basic

import "testing"

func TestGrid2D(t *testing.T) {
	rows, cols := 3, 3
	gen := NewGrid2DGenerator(rows, cols)
	g := gen.Generate()

	got := len(g.Nodes())
	expected := rows * cols
	if got != expected {
		t.Fatalf("Expected graph to have %d nodes, but got %d", expected, got)
	}

	got = len(g.Links())
	expected = 12
	if got != expected {
		t.Fatalf("Expected graph to have %d links, but got %d", expected, got)
	}
}
