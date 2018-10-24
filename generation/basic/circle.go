package basic

import (
	"github.com/divan/graphx/graph"
)

// CircleGenerator implements generator for circle graph.
type CircleGenerator struct {
	nodes int // number of nodes
}

// NewCircleGenerator creates new line generator for N nodes graph.
func NewCircleGenerator(n int) *CircleGenerator {
	return &CircleGenerator{
		nodes: n,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *CircleGenerator) Generate() *graph.Graph {
	data := graph.NewGraph()

	for i := 0; i < l.nodes+1; i++ {
		addNode(data, i)

		// close the circle
		if i == l.nodes {
			addLink(data, i-1, 0)
			continue
		}

		if i == 0 {
			continue
		}
		addLink(data, i-1, i)
	}

	return data
}
