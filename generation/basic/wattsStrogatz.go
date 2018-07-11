package basic

import (
	"math"
	"math/rand"

	"github.com/divan/graphx/graph"
)

// WattsStrogatzGenerator implements generator for Watts-Strogatz graph.
type WattsStrogatzGenerator struct {
	nodes              int // number of nodes
	conns              int // number of neigbours
	rewritePropability float64
}

// NewWattsStrogatzGenerator creates new Watts-Strogatz generator for N nodes graph.
func NewWattsStrogatzGenerator(n, conns int) *WattsStrogatzGenerator {
	if conns > n {
		panic("conns should be less then number of nodes")
	}
	return &WattsStrogatzGenerator{
		nodes:              n,
		conns:              conns,
		rewritePropability: 0.01,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *WattsStrogatzGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	for i := 0; i < l.nodes; i++ {
		addNode(g, i)

	}

	// connect each node conns/2 neigbors
	neigbors := int(math.Floor(float64(l.conns/2 + 1)))
	for i := 1; i < neigbors; i++ {
		for j := 0; j < l.nodes; j++ {
			to := int(math.Mod(float64(i+j), float64(l.nodes)))
			addLink(g, j, to)
		}
	}

	// rewire edges from each node
	neigbors = int(math.Floor(float64(l.conns/2 + 1)))
	for j := 1; j < neigbors; j++ {
		for i := 0; i < l.nodes; i++ {
			if rand.Float64() > l.rewritePropability {
				continue
			}

			from := i
			to := int(math.Mod(float64(i+j), float64(l.nodes)))
			newTo := rand.Intn(l.nodes)

			// TODO: switch to link indexes
			needsRewire := (newTo == i) || g.LinkExists(id(from), id(newTo))
			if needsRewire && (g.NodeLinks(id(from)) == l.nodes-1) {
				continue
			}

			for needsRewire {
				newTo = rand.Intn(l.nodes)
				needsRewire = (newTo == i) || g.LinkExists(id(from), id(newTo))
			}

			rewireLink(g, id(from), id(to), id(newTo))
		}
	}

	return g
}

// TODO: move it into graph package
func rewireLink(g *graph.Graph, from, to, newTo string) {
	links := g.Links()
	for i := range links {
		if links[i].From() == from && links[i].To() == to {
			links[i].Rewire(links[i].From(), newTo)
		} else if links[i].To() == from && links[i].From() == to {
			links[i].Rewire(newTo, links[i].To())
		}
	}
}
