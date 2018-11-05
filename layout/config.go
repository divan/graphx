package layout

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
	SpringStiffness: 0.05,
	SpringLen:       10.0,
	DragCoeff:       0.02,
}

// ForcesFromConfig constructs default forces from the given config.
func forcesFromConfig(config Config) []Force {
	repelling := NewGravityForce(config.Repelling, BarneHutMethod)
	springs := NewSpringForce(config.SpringStiffness, config.SpringLen, ForEachLink)
	drag := NewDragForce(config.DragCoeff, ForEachNode)
	return []Force{repelling, springs, drag}
}
