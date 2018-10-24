package graph

import "fmt"

// NodeHasLinks implements fast check if given node has any links.
func (g *Graph) NodeHasLinks(id string) bool {
	return g.nodeLinks[id] > 0
}

// NodeLinks returns number of links for node.
func (g *Graph) NodeLinks(id string) int {
	return g.nodeLinks[id]
}

// NodeIDHasLinks implements fast check if given node by ID has any links.
func (g *Graph) NodeIDHasLinks(id string) bool {
	return g.nodeLinks[id] > 0
}

// LinkExists returns true if there is a link between source and target.
func (g *Graph) LinkExists(from, to string) bool {
	for _, link := range g.links {
		if link.from == from && link.to == to ||
			link.to == from && link.from == to {
			return true
		}
	}
	return false
}

// LinkIndex returns link index by its source and target.
func (g *Graph) LinkIndex(from, to string) (int, error) {
	for i, link := range g.links {
		if link.from == from && link.to == to ||
			link.to == from && link.from == to {
			return i, nil
		}
	}

	return 0, fmt.Errorf("link %s->%s not found", from, to)
}

// LinkByIndices returns link index by its source and target indices.
func (g *Graph) LinkByIndices(from, to int) (int, error) {
	for i, link := range g.links {
		if link.fromIdx == from && link.toIdx == to ||
			link.toIdx == from && link.fromIdx == to {
			return i, nil
		}
	}

	return 0, fmt.Errorf("link %d->%d not found", from, to)
}

// NodeByID returns node index by its ID.
func (g *Graph) NodeByID(id string) (int, error) {
	for i, node := range g.nodes {
		if node.ID() == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("node %s not found", id)
}

// NodeIDByIdx returns node ID by its index.
func (g *Graph) NodeIDByIdx(idx int) (string, error) {
	if idx < 0 || idx > g.NumNodes()-1 {
		return "", fmt.Errorf("node for index %d not found", idx)
	}

	return g.nodes[idx].ID(), nil
}
