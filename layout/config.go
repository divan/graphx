package layout

import "github.com/divan/graphx/graph"

// Config specifies layout configuration for default set
// of forces.
type Config struct {
	Repelling                  float64
	SpringStiffness, SpringLen float64
	DragCoeff                  float64
}

// DefaultConfig sets reasonable values for small and medium graphs.
var DefaultConfig = Config{
	Repelling:       -10.0,
	SpringStiffness: 0.02,
	SpringLen:       10.0,
	DragCoeff:       0.8,
}

// NewFromConfig creates a new layout from the given config.
func NewFromConfig(g *graph.Graph, conf Config) *Layout {
	repelling := NewGravityForce(conf.Repelling, BarneHutMethod)
	springs := NewSpringForce(conf.SpringStiffness, conf.SpringLen, ForEachLink)
	drag := NewDragForce(conf.DragCoeff, ForEachNode)

	return New(g, repelling, springs, drag)
}
