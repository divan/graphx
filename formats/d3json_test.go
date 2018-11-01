package formats

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func BenchmarkImportJSON(b *testing.B) {
	files, err := readTestData()
	if err != nil {
		b.Fatalf("Failed to read testdata: %v", err)
	}

	for _, file := range files {
		b.Run(file, func(b *testing.B) {
			fd, err := os.Open(file)
			if err != nil {
				b.Fatal(err)
			}
			defer fd.Close()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				FromD3JSONReader(fd)
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
