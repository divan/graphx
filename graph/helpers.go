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
// TODO(divan): rename to NodeIdxByID
func (g *Graph) NodeByID(id string) (int, error) {
	idx, ok := g.nodeIdxByID[id]
	if !ok {
		// not in cache, attempt to find and cache
		for i := range g.nodes {
			if g.nodes[i].ID() == id {
				g.cacheNode(g.nodes[i], idx)
				return idx, nil
			}
		}
		return 0, fmt.Errorf("node %s not found", id)
	}
	return idx, nil
}

// NodeIDByIdx returns node ID by its index.
func (g *Graph) NodeIDByIdx(idx int) (string, error) {
	if idx < 0 || idx > g.NumNodes()-1 {
		return "", fmt.Errorf("node for index %d not found", idx)
	}

	return g.nodes[idx].ID(), nil
}

// Node returns Node by its string ID. It uses cache for faster lookup.
func (g *Graph) Node(id string) (Node, error) {
	idx, err := g.NodeByID(id)
	if err != nil {
		return nil, err
	}
	return g.nodes[idx], nil
}
