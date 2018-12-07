package layout

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/divan/graphx/generation/basic"
	"github.com/divan/graphx/graph"
)

func TestLayout(t *testing.T) {
	graph := basic.NewLineGenerator(2).Generate()
	repelling := NewGravityForce(-1.0, BarneHutMethod)
	//springs := NewSpringForce(0.01, 12.0, ForEachLink)
	l := NewWithForces(graph, repelling)
	l.objects["0"].SetPosition(0, 0, 0)
	l.objects["1"].SetPosition(100, 100, 100)

	for i := 0; i < 1000; i++ {
		l.UpdatePositions()
	}
}

func TestLayoutAdd(t *testing.T) {
	g := graph.NewGraph()
	g.AddNode(graph.NewBasicNode("node 0"))
	g.AddNode(graph.NewBasicNode("node 1"))
	err := g.AddLink("node 0", "node 1")
	if err != nil {
		t.Fatalf("Add link failed: %v", err)
	}
	repelling := NewGravityForce(-1.0, ForEachNode)
	l := NewWithForces(g, repelling)

	if len(l.objects) != 2 {
		t.Fatalf("objects map expected to be of %d length, but is of %d", 2, len(l.objects))
	}
	if len(l.keys) != 2 {
		t.Fatalf("keys expected to be of %d length, but is of %d", 2, len(l.objects))
	}
	if len(l.links) != 1 {
		t.Fatalf("expected to have %d links, but has %d", 1, len(l.links))
	}

	l.AddNode(graph.NewBasicNode("node 2"))
	if len(l.objects) != 3 {
		t.Fatalf("objects map expected to be of %d length, but is of %d", 3, len(l.objects))
	}
	if len(l.keys) != 3 {
		t.Fatalf("keys expected to be of %d length, but is of %d", 3, len(l.objects))
	}
}

func BenchmarkUpdatePositions(b *testing.B) {
	files, err := readTestData()
	if err != nil {
		b.Fatalf("Failed to read testdata: %v", err)
	}

	for _, file := range files {
		g, err := FromD3JSON(file)
		if err != nil {
			b.Fatalf("Failed to build graph: %v", err)
		}
		l := New(g, DefaultConfig)
		b.Run(file, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				l.UpdatePositions()
			}
		})
	}
}

// readTestData returns filenames (with relative path)
// for files in testdata/ directory.
func readTestData() ([]string, error) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		return nil, fmt.Errorf("readdir: %s", err)
	}

	var ret []string
	for _, file := range files {
		path := filepath.Join("testdata", file.Name())
		ret = append(ret, path)
	}
	return ret, nil
}

// TODO(divan): remove this from here, once package formats is refactored,
// so there is no circular dependency.

// FromD3JSON creates a graph from the given JSON file.
// It recognizes simple JSON structure suitable for D3 examples,
// basically just `id`, `group` and `weight` fields for nodes.
func FromD3JSON(file string) (*graph.Graph, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close() //nolint: errcheck

	return FromD3JSONReader(fd)
}

// FromD3JSONReader creates a graph from the given JSON file.
func FromD3JSONReader(r io.Reader) (*graph.Graph, error) {
	// decode into temporary struct to process
	var res struct {
		Nodes []*graph.BasicNode `json:"nodes"`
		Links []*struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"links"`
	}
	err := json.NewDecoder(r).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("decode JSON: %v", err)
	}

	if len(res.Nodes) == 0 {
		return nil, errors.New("empty graph")
	}

	// convert links IDs into indices
	g := graph.NewGraphMN(len(res.Nodes), len(res.Links))

	for _, node := range res.Nodes {
		g.AddNode(node)
	}

	for _, link := range res.Links {
		err := g.AddLink(link.Source, link.Target)
		if err != nil {
			return nil, fmt.Errorf("add link: %v", err)
		}
	}

	return g, err
}
