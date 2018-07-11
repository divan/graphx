package net

import "github.com/divan/graphx/graph"

// SplitBrainGenerator implements a network with
// two separated parts (aka hemispheres) and only
// one link between them (can be cut down manually
// in a resulting JSON if needed).
type SplitBrainGenerator struct {
	hosts       int
	connections int
	startIP     string
}

// NewSplitBrainGenerator creats new split-brain network generator, and applies
// given paramteres to the both parts of the network.
func NewSplitBrainGenerator(hosts, conns int, startIP string) *SplitBrainGenerator {
	return &SplitBrainGenerator{
		hosts:       hosts,
		connections: conns,
		startIP:     startIP,
	}
}

// Generate generates (almost) split-brain network with known number of
// hosts with known number of connections each. Implements Generator
// interface.
func (d *SplitBrainGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	gen := NewIPGenerator(d.startIP)

	// generate hosts for left subnetwork
	for i := 0; i < d.hosts/2; i++ {
		ip := gen.NextAddress()
		node := NewNode(ip)
		g.AddNode(node)
	}

	// generate links
	for i := 0; i < d.hosts; i++ {
		for j := 0; j < d.connections; j++ {
			idx, err := nextIdx(g, i, 0, d.hosts)
			if err == nil {
				addLink(g, i, idx)
			}
		}
	}

	// right subnetwork
	for i := 0; i < d.hosts/2; i++ {
		ip := gen.NextAddress()
		node := NewNode(ip)
		g.AddNode(node)
	}
	for i := d.hosts / 2; i < d.hosts; i++ {
		for j := 0; j < d.connections; j++ {
			idx, err := nextIdx(g, i, d.hosts/2, d.hosts)
			if err == nil {
				addLink(g, i, idx)
			}
		}
	}

	return g
}
