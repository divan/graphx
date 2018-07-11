package formats

import (
	"github.com/divan/graphx/graph"
	"github.com/divan/graphx/layout"
)

// GraphExporter defines exporters for graph.
type GraphExporter interface {
	ExportGraph(*graph.Graph) error
}

// LayoutExporter defines exporters for layout.
type LayoutExporter interface {
	ExportLayout(layout.Layout) error
}
