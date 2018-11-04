package graph

func (g *Graph) ResetCache() {
	g.nodeLinks = make(map[string]int)
	g.nodeIdxByID = make(map[string]int)
}

func (g *Graph) cacheNode(node Node, idx int) {
	g.nodeIdxByID[node.ID()] = idx
}

func (g *Graph) cacheLink(from, to string) {
	g.nodeLinks[from]++
	g.nodeLinks[to]++
}
