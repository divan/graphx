package graph

// Link represents single link between two nodes.
type Link struct {
	from string
	to   string

	fromIdx int
	toIdx   int
}

// NewLink constructs new Link object.
// Note, this function doesn't know actual nodes, so it doesn't
// check for indices validity.
func NewLink(from, to string) *Link {
	return &Link{
		from: from,
		to:   to,
	}
}

// AddLink adds new link to the graph and validates input
// indices.
func (g *Graph) AddLink(from, to string) error {
	// TODO: add node if ID is unexistent

	link := NewLink(from, to)

	var err error
	link.fromIdx, err = g.NodeByID(from)
	if err != nil {
		return err
	}
	link.toIdx, err = g.NodeByID(to)
	if err != nil {
		return err
	}

	g.links = append(g.links, link)
	return nil
}

// From returns link's source ID.
func (l *Link) From() string { return l.from }

// To returns link's target ID.
func (l *Link) To() string { return l.to }

// FromIdx returns link's source index.
func (l *Link) FromIdx() int { return l.fromIdx }

// ToIdx returns link's target index.
func (l *Link) ToIdx() int { return l.toIdx }

// Rewire allows explicitly change edge.
func (l *Link) Rewire(from, to string) {
	l.from = from
	l.to = to
}
