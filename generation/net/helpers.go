package net

import (
	"fmt"

	"github.com/divan/graphx/graph"
)

func addNode(g *graph.Graph, i int) {
	node := graph.NewBasicNode(id(i))
	g.AddNode(node)
}

func addLink(g *graph.Graph, i, j int) {
	g.AddLink(id(i), id(j))
}

func id(i int) string {
	return fmt.Sprintf("%d", i)
}
