package layout

import (
	"fmt"

	"github.com/divan/graphx/graph"
)

// ForceRule define algorithm/rules to apply force on a graph. Force can be applied an a variety of different ways and this abstraction should ideally catch and encapsulate all these differences.
//
// vectors and debugInfo are passed for optimization purposes, to avoid allocating new memory.
type ForceRule func(
	force Force,
	objects map[string]*Object,
	links []*graph.Link)

// ForEachLink applies force to both ends of each link in the graph, with positive and negative signs respectively.
var ForEachLink = func(
	force Force,
	objects map[string]*Object,
	links []*graph.Link,
	vectors map[int]*ForceVector,
	debugInfo ForcesDebugData) {
	for _, link := range links {
		idFrom := link.From()
		idTo := link.To()

		from := objects[idFrom]
		to := objects[idTo]
		f := force.Apply(from, to)

		// Update force vectors
		from.force.Add(f)
		to.force.Sub(f)

		// Update debug information
		/*
			name := force.Name()
			debugInfo.Append(link.From, name, *f)
			debugInfo.Append(link.To, name, f.Negative())
		*/
	}
}

// BarneHutMethod applies force for each node agains each node,
// using Barne-Hut optimization method.
var BarneHutMethod = func(
	force Force,
	objects map[string]*Object,
	links []*graph.Link) {

	otree := NewOctreeFromNodes(objects, force)

	for i, node := range objects {
		f, err := otree.CalcForce(i)
		if err != nil {
			fmt.Println("[ERROR] Force calc failed:", i, err)
			continue
		}

		// Update force vectors
		f1 := objects[node.ID].force
		objects[node.ID].force = f1.Add(f)

		// Update debug information
		//name := force.Name()
		//debugInfo.Append(node.Idx, name, *f)
	}
}

// ForEachNode applies force to every node in the graph.
var ForEachNode = func(
	force Force,
	objects map[string]*Object,
	links []*graph.Link,
	vectors map[int]*ForceVector,
	debugInfo ForcesDebugData) {
	for id, node := range objects {
		f := force.Apply(node, nil)

		// Update force vectors
		ff := objects[id].force
		objects[id].force = ff.Add(f)

		// Update debug information
		//name := force.Name()
		//debugInfo.Append(i, name, *f)
	}
}
