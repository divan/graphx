package graph

// Graph represents graph data.
type Graph struct {
	nodes []Node
	links []*Link

	nodeLinks map[string]int
}

// NewGraph creates empty graph data.
func NewGraph() *Graph {
	return &Graph{
		nodeLinks: make(map[string]int),
	}
}

// NewGraphMN creates empty graph for M nodes and N links.
// It preallocates memory for the specified sizes.
func NewGraphMN(m, n int) *Graph {
	return &Graph{
		nodes:     make([]Node, 0, m),
		links:     make([]*Link, 0, n),
		nodeLinks: make(map[string]int),
	}
}

// Nodes returns graph nodes
func (g *Graph) Nodes() []Node {
	return g.nodes
}

// Links returns graph links.
func (g *Graph) Links() []*Link {
	return g.links
}

// UpdateCache runs various optimization-related
// calculations, caching etc.
func (g *Graph) UpdateCache() {
	g.nodeLinks = make(map[string]int)
	for _, link := range g.links {
		g.nodeLinks[link.from]++
		g.nodeLinks[link.to]++
	}
}
