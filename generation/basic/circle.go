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

	for i := 0; i < l.nodes; i++ {
		addNode(data, i)

		j := i + 1
		if i == l.nodes-1 {
			j = 0
		}
		data.AddLink(id(i), id(j))
	}

	return data
}
