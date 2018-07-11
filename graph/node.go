package graph

// Node defines the graph node. Any type implementing this
// interface can be used as a graph node.
type Node interface {
	ID() string
}

// GroupedNode represents node that have 'group' attribute.
type GroupedNode interface {
	Group() int
}

// WeightedNode represents node that have 'weight' attribute.
type WeightedNode interface {
	Weight() int
}

// AddNode adds new node to graph.
func (g *Graph) AddNode(node Node) {
	g.nodes = append(g.nodes, node)
}

// AddNodes adds new nodes to graph.
func (g *Graph) AddNodes(nodes ...Node) {
	g.nodes = append(g.nodes, nodes...)
}

// BasicNode represents basic built-in node type for simple cases.
type BasicNode struct {
	ID_     string `json:"id"`
	Group_  int    `json:"group,omitempty"`
	Weight_ int    `json:"weight,omitempty"`
}

// ID implements Node for BasicNode.
func (b *BasicNode) ID() string { return b.ID_ }

// Group implements GroupNode for BasicNode.
func (b *BasicNode) Group() int { return b.Group_ }

// Weight implements WeightedNode for BasicNode.
func (b *BasicNode) Weight() int { return b.Weight_ }

// NewBasicNode creaates a new basic node with given ID.
func NewBasicNode(id string) *BasicNode {
	return &BasicNode{
		ID_: id,
	}
}
