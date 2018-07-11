package formats

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/divan/graphx/graph"
)

func TestD3JSONExport(t *testing.T) {
	var buf bytes.Buffer
	w := NewD3JSON(&buf, false)
	g := testGraph()
	err := w.ExportGraph(g)
	if err != nil {
		t.Fatalf("Exporting graph to D3 JSON failed: %v", err)
	}

	expected := 111
	if buf.Len() != expected {
		t.Log(buf.String())
		t.Fatalf("Exported JSON should of length %d, but got: %d", expected, buf.Len())
	}

	g1, err := FromD3JSONReader(&buf)
	if err != nil {
		t.Fatalf("Importing graph from D3 JSON failed: %v", err)
	}

	if len(g1.Links()) != len(g.Links()) {
		t.Fatalf("Imported and exported graphs has different amount of links")
	}
	if len(g1.Nodes()) != len(g.Nodes()) {
		t.Fatalf("Imported and exported graphs has different amount of nodes")
	}
}

func testGraph() *graph.Graph {
	g := graph.NewGraph()
	g.AddNode(node(1))
	g.AddNode(node(2))
	g.AddNode(node(3))

	g.AddLink("1", "2")
	g.AddLink("2", "3")
	return g
}

func node(i int) *graph.BasicNode { return &graph.BasicNode{ID_: fmt.Sprintf("%d", i)} }

func TestNewJSONGraph(t *testing.T) {
	buf := bytes.NewBufferString(`{
		"nodes": [ {"id": "A", "weight": 10}, {"id": "B"}, {"id": "C"}, {"id": "D"} ],
		"links": [ {"source": "A", "target": "B"}, {"source": "C", "target": "D"}, {"source": "C", "target": "A"}]
	}`)
	graph, err := FromD3JSONReader(buf)
	if err != nil {
		t.Fatal(err)
	}

	nodes := graph.Nodes()
	if len(nodes) != 4 {
		t.Fatalf("Expect graph to have %d nodes, but got %d", 4, len(nodes))
	}

	links := graph.Links()
	if len(links) != 3 {
		t.Fatalf("Expect graph to have %d links, but got %d", 3, len(links))
	}

	linksCounter := map[string]int{
		"A": 2,
		"B": 1,
		"C": 2,
		"D": 1,
	}
	for _, node := range nodes {
		got := graph.NodeLinks(node.ID())
		expected := linksCounter[node.ID()]
		if got != expected {
			t.Fatalf("Expected number of links to be %d, but got %d for node '%s'",
				expected, got, node.ID())
		}
	}
}

func TestNewJSONGraphLarge(t *testing.T) {
	N := int(10e4)
	r := generateLargeGraphJSON(N, 2*N)

	graph, err := FromD3JSONReader(r)
	if err != nil {
		t.Fatal(err)
	}

	nodes := graph.Nodes()
	if len(nodes) != N {
		t.Fatalf("Expect graph to have %v nodes, but got %d", N, len(nodes))
	}

	links := graph.Links()
	if len(links) != 2*N {
		t.Fatalf("Expect graph to have %v links, but got %d", 2*N, len(links))
	}
}

func BenchmarkImportJSONGraph(b *testing.B) {
	N := int(10e4)
	r := generateLargeGraphJSON(N, 2*N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FromD3JSONReader(r)
	}
}

// generate large graph with given number of nodes and links.
func generateLargeGraphJSON(nodes, links int) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString(`{"nodes":[`)
	for i := 0; i < nodes; i++ {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"id": "same"}`)
	}
	buf.WriteString(`],"links":[`)
	for i := 0; i < links; i++ {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"source": "same", "target": "same"}`)
	}
	buf.WriteString(`]}`)
	return buf
}
