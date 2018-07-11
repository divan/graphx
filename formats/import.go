package formats

import "github.com/divan/graphx/graph"

// GraphImporter defines importers for graph.
type GraphImporter interface {
	ImportGraph() (*graph.Graph, error)
}
