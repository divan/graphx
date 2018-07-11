package layout

// ForcesDebugData holds debug information for forces
// for all nodes (map key is a node index).
type ForcesDebugData map[int][]*ForceDebugInfo

// ForceDebugInfo contains information for debugging calculated forces.
type ForceDebugInfo struct {
	Name string `json:"name"`
	ForceVector
}

// Append appends new ForceDebugInfo to the debug data.
func (f ForcesDebugData) Append(idx int, name string, force ForceVector) {
	if f == nil {
		f = make(map[int][]*ForceDebugInfo)
	}

	debugInfo := &ForceDebugInfo{
		Name:        name,
		ForceVector: force,
	}

	f[idx] = append(f[idx], debugInfo)
}
