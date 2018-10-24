package basic

import (
	"github.com/divan/graphx/graph"
)

// LineGenerator implements generator for line graph.
type LineGenerator struct {
	nodes int // number of nodes
}

// NewLineGenerator creates new line generator for N nodes graph.
func NewLineGenerator(n int) *LineGenerator {
	return &LineGenerator{
		nodes: n,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *LineGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	for i := 0; i < l.nodes; i++ {
		addNode(g, i)

		if i == 0 {
			continue
		}

		addLink(g, i-1, i)
	}

	return g
}
