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
	links []*graph.Link) {
	for _, link := range links {
		idFrom := link.From()
		idTo := link.To()

		from := objects[idFrom]
		to := objects[idTo]
		f := force.Apply(from, to)

		// Update force vectors
		objects[idFrom].force.Add(f)
		objects[idTo].force.Sub(f)
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
			break
		}

		objects[node.ID].force.Add(f)
	}
}

// ForEachNode applies force to every node in the graph.
var ForEachNode = func(
	force Force,
	objects map[string]*Object,
	links []*graph.Link) {
	for id, node := range objects {
		f := force.Apply(node, nil)

		objects[id].force.Add(f)
	}
}

// EachOnEach applies every node force to every other node in the graph.
// It's slow, and added just to compare results with more optimized versions.
var EachOnEach = func(
	force Force,
	objects map[string]*Object,
	links []*graph.Link) {
	var newObjects = make(map[string]*Object)
	for id, node := range objects {
		newObjects[id] = NewObject(0, 0, 0)
		newObjects[id].force.DX = node.force.DX
		newObjects[id].force.DY = node.force.DY
		newObjects[id].force.DZ = node.force.DZ
	}

	for id1 := range objects {
		for id2 := range objects {
			if id1 == id2 {
				continue
			}

			f := force.Apply(objects[id1], objects[id2])

			// Update force vectors
			newObjects[id1].force.Add(f)
			newObjects[id2].force.Sub(f)
		}
	}

	for id, node := range newObjects {
		objects[id].force.DX = node.force.DX
		objects[id].force.DY = node.force.DY
		objects[id].force.DZ = node.force.DZ
	}
}
