package basic

import (
	"github.com/divan/graphx/graph"
)

// KingGenerator implements generator for king graph graph.
// https://en.wikipedia.org/wiki/King%27s_graph
type KingGenerator struct {
	rows int
	cols int
}

// NewKingGenerator creates new king graph generator for known number of rows and cols.
func NewKingGenerator(rows, cols int) *KingGenerator {
	return &KingGenerator{
		rows: rows,
		cols: cols,
	}
}

// NewKingGeneratorN creates new graph generator for N nodes.
func NewKingGeneratorN(n int) *KingGenerator {
	rows, cols := estimateRowsCols(n)
	return &KingGenerator{
		rows: rows,
		cols: cols,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *KingGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	for r := 0; r < l.rows; r++ {
		for c := 0; c < l.cols; c++ {
			idx := c + r*l.rows

			addNode(g, idx)

			// same row, left col
			if c > 0 {
				addLink(g, idx, c-1+r*l.rows)
			}

			// prev col, prev row
			if c > 0 && r > 0 {
				addLink(g, idx, c-1+(r-1)*l.rows)
			}
			// next col, prev row
			if c < l.cols-1 && r > 0 {
				addLink(g, idx, c+1+(r-1)*l.rows)
			}

			// same col, prev row
			if r > 0 {
				addLink(g, idx, c+(r-1)*l.rows)
			}
		}
	}

	return g
}
