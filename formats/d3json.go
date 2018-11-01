package formats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/divan/graphx/graph"
)

// D3JSON implements GraphExporter for D3 JSON format. It's a simple JSON, without any specification,
// just used by a bunch of JS projects like D3 and ngraph.
//
// See https://bl.ocks.org/mbostock/4062045 for example.
type D3JSON struct {
	writer   io.Writer
	indented bool
}

// ToD3JSON is a helper for D3JSON exporter for saving graph into the D3 JSON formating to
// the given file.
func ToD3JSON(g *graph.Graph, file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer func(fd *os.File) {
		err = fd.Close()
		if err != nil {
			log.Println(fmt.Errorf("close file: %s", err))
		}
	}(fd)

	d3 := NewD3JSON(fd, true) // indent enabled
	return d3.ExportGraph(g)
}

// NewD3JSON creates new D3JSON exporter. Indented specifies if produced JSON should be indented.
func NewD3JSON(w io.Writer, indented bool) *D3JSON {
	return &D3JSON{
		writer:   w,
		indented: indented,
	}
}

// ExportGraph converts graph into D3 JSON format. Implements GraphExporter interface.
func (d *D3JSON) ExportGraph(g *graph.Graph) error {
	type link struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}
	var data struct {
		Nodes []graph.Node `json:"nodes"`
		Links []*link      `json:"links"`
	}

	data.Nodes = g.Nodes()
	data.Links = make([]*link, g.NumLinks())
	for i, l := range g.Links() {
		data.Links[i] = &link{
			Source: l.From(),
			Target: l.To(),
		}
	}

	enc := json.NewEncoder(d.writer)
	if d.indented {
		enc.SetIndent("", "  ")
	}

	return enc.Encode(data)
}

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
		return nil, err
	}

	if len(res.Nodes) == 0 {
		return nil, errors.New("empty graph")
	}

	// convert links IDs into indices
	g := graph.NewGraphMN(len(res.Nodes), len(res.Links))

	for _, node := range res.Nodes {
		g.AddNode(node)
	}

	g.UpdateCache()

	for _, link := range res.Links {
		err := g.AddLink(link.Source, link.Target)
		if err != nil {
			return nil, err
		}
	}

	g.UpdateCache()

	return g, err
}
