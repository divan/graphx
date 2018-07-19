package layout

import (
	"fmt"
	"math"

	"github.com/divan/graphx/graph"
)

// NewAuto will init 3D layout and automatically estimate forces and it's paramteres
// for this particular graph.
func NewAuto(g *graph.Graph) *Layout {
	worldSize := float64(2000) // TODO: this should be synced/communicated to with frontend somehow

	graphWidth := estimateGraphWidth(g)

	optimalEdge := estimateOptimalEdge(worldSize, len(g.Links()))

	repForce := -(worldSize / graphWidth / 40)
	fmt.Println("Optimal edge:", optimalEdge)
	fmt.Println("Graph width (not real):", graphWidth)
	fmt.Println("Repelling force:", repForce)
	repelling := NewGravityForce(repForce, BarneHutMethod)
	springs := NewSpringForce(0.02, optimalEdge, ForEachLink)
	drag := NewDragForce(0.8, ForEachNode)

	return New(g, repelling, springs, drag)
}

func estimateOptimalEdge(width float64, links int) float64 {
	k := 0.1
	return k * math.Sqrt((width*width)/float64(links))
}

func estimateGraphWidth(g *graph.Graph) float64 {
	// TODO: implement it via "longest short path"
	// use cube root for now
	return math.Cbrt(float64(len(g.Nodes())))
}
