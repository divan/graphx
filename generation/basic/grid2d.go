package basic

import (
	"math"

	"github.com/divan/graphx/graph"
)

// Grid2DGenerator implements generator for 2D grid graph.
type Grid2DGenerator struct {
	rows int
	cols int
}

// NewGrid2DGenerator creates new grid graph generator for known number of rows and cols.
func NewGrid2DGenerator(rows, cols int) *Grid2DGenerator {
	return &Grid2DGenerator{
		rows: rows,
		cols: cols,
	}
}

// NewGrid2DGeneratorN creates new graph generator for N nodes.
func NewGrid2DGeneratorN(n int) *Grid2DGenerator {
	rows, cols := estimateRowsCols(n)
	return &Grid2DGenerator{
		rows: rows,
		cols: cols,
	}
}

// estimateRowsCols tries to find multiplies for n closest to square.
// TODO: make it efficient and correct :)
func estimateRowsCols(n int) (int, int) {
	root := math.Round(math.Sqrt(float64(n)))
	if root < 2 {
		return 1, 1
	}

	res := math.Mod(root, float64(n))
	if math.Round(res) < 2 {
		return int(root), 1
	}

	return int(root), int(res)
}

// Generate generates the data for graph. Implements Generator interface.
func (l *Grid2DGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	for y := 0; y < l.rows; y++ {
		for x := 0; x < l.cols; x++ {
			idx := x + y*l.cols
			addNode(g, idx)

			if x > 0 {
				leftIdx := x - 1 + y*l.cols
				addLink(g, idx, leftIdx)
			}
			if y > 0 {
				topIdx := x + (y-1)*l.cols
				addLink(g, idx, topIdx)
			}
		}
	}

	return g
}
