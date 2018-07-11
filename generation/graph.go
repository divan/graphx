package generation

import "github.com/divan/graphx/graph"

// GraphGenerator represents generator that generates
// graph data with nodes and links.
type GraphGenerator interface {
	Generate() *graph.Graph
}
