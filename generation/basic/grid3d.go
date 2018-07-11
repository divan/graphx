package basic

import (
	"math"

	"github.com/divan/graphx/graph"
)

// Grid3DGenerator implements generator for 3D grid graph.
type Grid3DGenerator struct {
	rows   int
	cols   int
	levels int
}

// NewGrid3DGenerator creates new grid graph generator for known number of rows, cols and levels.
func NewGrid3DGenerator(rows, cols, levels int) *Grid3DGenerator {
	return &Grid3DGenerator{
		rows:   rows,
		cols:   cols,
		levels: levels,
	}
}

// NewGrid3DGeneratorN creates new grid 3D graph generator for N nodes.
func NewGrid3DGeneratorN(n int) *Grid3DGenerator {
	rows, cols, levels := estimateRowsColsLevels(n)
	return &Grid3DGenerator{
		rows:   rows,
		cols:   cols,
		levels: levels,
	}
}

// estimateRowsCols tries to find multiplies for n closest to cube.
// TODO: make it efficient and correct :)
func estimateRowsColsLevels(n int) (int, int, int) {
	root := math.Round(math.Cbrt(float64(n)))
	if root < 2 {
		return 1, 1, 1
	}

	return int(root), int(root), int(root)
}

// Generate generates the data for graph. Implements Generator interface.
func (l *Grid3DGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	for k := 0; k < l.levels; k++ {
		for i := 0; i < l.rows; i++ {
			for j := 0; j < l.cols; j++ {
				level := k * l.rows * l.cols
				idx := j + i*l.rows + level
				addNode(g, idx)

				if j > 0 {
					addLink(g, idx, j-1+i*l.rows+level)
				}
				if i > 0 {
					addLink(g, idx, j+(i-1)*l.rows+level)
				}
				if k > 0 {
					addLink(g, idx, j+i*l.rows+(k-1)*l.rows*l.cols)
				}
			}
		}
	}

	return g
}
