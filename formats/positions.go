package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/divan/graphx/layout"
)

// FromPositionsJSON reads points positions from io.Reader. Index should correspond the original graph
// nodes indicies.
func FromPositionsJSON(r io.Reader) ([]*layout.Point, error) {
	var ret []*layout.Point
	err := json.NewDecoder(r).Decode(&ret)
	return ret, err
}

// FromPositionsJSONFile reads points positions from file. Index should correspond the original graph
// nodes indicies.
func FromPositionsJSONFile(file string) ([]*layout.Point, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open positions json: %v", err)
	}
	defer fd.Close()
	return FromPositionsJSON(fd)
}

// ToPositionsJSON writes points positions to io.Writer. Index should correspond the original graph
// nodes indicies.
func ToPositionsJSON(positions []*layout.Point, w io.Writer) error {
	return json.NewEncoder(w).Encode(positions)
}

// ToPositionsJSONFile writes points positions to the file. Index should correspond the original graph
// nodes indicies.
func ToPositionsJSONFile(positions []*layout.Point, file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create positions json: %v", err)
	}
	return ToPositionsJSON(positions, fd)
}
