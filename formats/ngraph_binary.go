package formats

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/divan/graphx/graph"
	"github.com/divan/graphx/layout"
)

// NgraphBinary stores graph data as binary files, compatible
// with anvaka/ngraph library: positions.bin, links.bin, labels.json and meta.json.
//
// Files are stored as binary data with little-endian encoding, to minimize size and
// optimize for large graphs.
//
// See https://github.com/anvaka/ngraph.offline.layout for more information
type NgraphBinary struct {
	Dir string // output directory for .bin files
}

// NewNgraphBinary creates new ngraph binary formatter, and sets output dir to dir.
// Directory must exist.
func NewNgraphBinary(dir string) (*NgraphBinary, error) {
	fs, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("check output dir: %v", err)
	}
	if !fs.IsDir() {
		log.Fatalf("'%s' is not a dir, aborting...", dir)
		return nil, fmt.Errorf("'%s' is not a dir", dir)
	}

	return &NgraphBinary{
		Dir: dir,
	}, nil
}

// ExportGraph exports graph structure and data into binary files (links.bin and labels.json) in
// an output directory.
// Implements GraphExporter.
func (n *NgraphBinary) ExportGraph(g *graph.Graph) error {
	err := n.writeLinksBin(g)
	if err != nil {
		return err
	}

	return n.writeLabels(g)
}

// ExportLayout writes position data into 'positions.bin' file in the
// following way: XYZXYZXYZ... where X, Y and Z are coordinates
// for each node in signed 32 bit integer Little Endian format.
// Implements export.LayoutExporter interface.
func (n *NgraphBinary) ExportLayout(l layout.Layout) error {
	file := filepath.Join(n.Dir, "positions.bin")
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	iw := newInt32LEWriter(fd)

	nodes := l.Positions()
	for k := range nodes {
		iw.Write(int32(nodes[k].X))
		iw.Write(int32(nodes[k].Y))
		iw.Write(int32(nodes[k].Z))
		if iw.err != nil {
			return err
		}
	}

	return nil
}

// writeLinksBin writes links information into `links.bin` file in the
// following way: Sidx,L1idx,L2idx,S2idx,L1idx... where SNidx - is the
// start node index, and LNidx - is the other link end node index.
func (n *NgraphBinary) writeLinksBin(g *graph.Graph) error {
	file := filepath.Join(n.Dir, "links.bin")
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	iw := newInt32LEWriter(fd)
	_ = iw
	for _, node := range g.Nodes() {
		if !g.NodeHasLinks(node.ID()) {
			continue
		}

		/* TODO: implement links cache
		iw.Write(int32(-(i + 1)))
		for _, link := range data.Links() {
			if link.From() == i {
				iw.Write(int32(link.To() + 1))
			}
		}
		if iw.err != nil {
			return err
		}
		*/
	}
	return nil
}

// writeLabels writes node ids (labels) information into `labels.json` file
// as an array of strings.
func (n *NgraphBinary) writeLabels(g *graph.Graph) error {
	file := filepath.Join(n.Dir, "labels.json")
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	var labels []string
	for i := range g.Nodes() {
		labels = append(labels, g.Nodes()[i].ID())
	}
	return json.NewEncoder(fd).Encode(labels)
}

// int32LEWriter implements binary writer for signed little-endian 32bit integers.
// It's used for ngraph_binary format.
type int32LEWriter struct {
	w   io.Writer
	err error
}

// newInt32LEWriter creates new inte32LEWriter.
func newInt32LEWriter(w io.Writer) *int32LEWriter {
	return &int32LEWriter{
		w: w,
	}
}

// Write writes given int32 into writer in binary format.
func (iw *int32LEWriter) Write(number int32) {
	if iw.err != nil {
		return
	}

	err := binary.Write(iw.w, binary.LittleEndian, number)
	if err != nil {
		iw.err = err
	}

	return
}
