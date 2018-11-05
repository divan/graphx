package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/divan/graphx/layout"
)

// Position represents only X, Y, Z coordinates of object.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// FromPositionsJSON reads points positions from io.Reader. Index should correspond the original graph
// nodes indicies.
func FromPositionsJSON(r io.Reader) ([]Position, error) {
	var ret []Position
	err := json.NewDecoder(r).Decode(&ret)
	return ret, err
}

// FromPositionsJSONFile reads points positions from file. Index should correspond the original graph
// nodes indicies.
func FromPositionsJSONFile(file string) ([]Position, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open positions json: %v", err)
	}
	defer fd.Close()
	return FromPositionsJSON(fd)
}

// ToPositionsJSON writes points positions to io.Writer. Index should correspond the original graph
// nodes indicies.
func ToPositionsJSON(positions []*layout.Position, w io.Writer) error {
	return json.NewEncoder(w).Encode(positions)
}

// ToPositionsJSONFile writes points positions to the file. Index should correspond the original graph
// nodes indicies.
func ToPositionsJSONFile(positions []*layout.Position, file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create positions json: %v", err)
	}
	defer fd.Close()
	return ToPositionsJSON(positions, fd)
}

// ToPositionsNGraph writes points positions to the io.Writer in the NGraph binary format.
func ToPositionsNGraph(positions []*layout.Position, w io.Writer) error {
	iw := newInt32LEWriter(w)

	for k := range positions {
		iw.Write(int32(positions[k].X))
		iw.Write(int32(positions[k].Y))
		iw.Write(int32(positions[k].Z))
		if iw.err != nil {
			return fmt.Errorf("write Int32LE: %v", iw.err)
		}
	}

	return nil
}

// ToPositionsNGraphFile writes points positions to the file in the NGraph binary format.
func ToPositionsNGraphFile(positions []*layout.Position, file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create positions ngraph binary: %v", err)
	}
	defer fd.Close()
	return ToPositionsNGraph(positions, fd)
}
