package layout

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/divan/graphx/graph"
	"gopkg.in/cheggaaa/pb.v1"
)

// stableThreshold determines the movement diff needed to
// call the system stable
const stableThreshold = 2.001

// Layout implements Layout interface for force-directed 3D graph.
type Layout struct {
	g *graph.Graph

	objects map[string]*Object // node ID as a key
	links   []*graph.Link
	forces  []Force
}

// New initializes 3D layout with objects data and set of forces.
func New(g *graph.Graph, forces ...Force) *Layout {
	l := &Layout{
		g:       g,
		objects: make(map[string]*Object),
		links:   g.Links(),
		forces:  forces,
	}

	l.initPositions()

	return l
}

// initPositions inits layout graph from the original graph data.
func (l *Layout) initPositions() {
	for _, node := range l.g.Nodes() {
		l.AddNode(node)
	}

	l.resetForces()
}

// AddNode handles adding new node to the existing layout.
func (l *Layout) AddNode(node graph.Node) error {
	if _, err := l.g.NodeByID(node.ID()); err == nil {
		return errors.New("node exists")
	}
	l.g.AddNode(node)

	l.addObject(node)
	return nil
}

// TODO: add link/remove link

func (l *Layout) addObject(node graph.Node) {
	lastIdx := len(l.objects) // use last item index for calculating pseudo-random positions

	// TODO: handle weight

	x, y, z := randomPosition(lastIdx)

	o := NewObjectID(x, y, z, node.ID())

	l.objects[node.ID()] = o
}

// randomPosition generates x,y,z coordinates pseudo-randomly spread around the
// spherical surface.
func randomPosition(i int) (x, y, z int) {
	radius := 10 * math.Cbrt(float64(i))
	rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
	yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

	x = int(radius * math.Cos(rollAngle))
	y = int(radius * math.Sin(rollAngle))
	z = int(radius * math.Sin(yawAngle))
	return
}

// Calculate runs positions' recalculations iteratively until the
// system minimizes it's energy.
func (l *Layout) Calculate() {
	// tx is the total movement, which should drop to the minimum
	// at the minimal energy state
	fmt.Println("Simulation started...")
	var (
		now    = time.Now()
		count  int
		prevTx float64
	)
	for tx := math.MaxFloat64; math.Abs(tx-prevTx) >= stableThreshold; {
		prevTx = tx
		tx = l.UpdatePositions()
		log.Println("PrevTx, tx:", tx, ", diff:", math.Abs(tx-prevTx))
		count++
		if count%1000 == 0 {
			since := time.Since(now)
			fmt.Printf("Iterations: %d, tx: %f, time: %v\n", count, tx, since)
		}
	}
	fmt.Printf("Simulation finished in %v, run %d iterations\n", time.Since(now), count)
}

// CalculateN run positions' recalculations exactly N times.
func (l *Layout) CalculateN(n int) {
	fmt.Println("Simulation started...")
	bar := pb.StartNew(n)
	for i := 0; i < n; i++ {
		l.UpdatePositions()
		bar.Increment()
	}
	bar.FinishPrint("Simulation finished")

}

// UpdatePositions recalculates nodes' positions, applying all the forces.
// It returns average amount of movement generated by this step.
func (l *Layout) UpdatePositions() float64 {
	//l.resetForces()

	for _, force := range l.forces {
		apply := force.Rule()
		apply(force, l.objects, l.links)
	}

	for _, v := range l.objects {
		fmt.Println("Force", v.force)
		break
	}

	return l.integrate()
}

func (l *Layout) resetForces() {
	for k := range l.objects {
		l.objects[k].force = ZeroForce
	}
}

// AddForce adds force to the internal list of forces.
func (l *Layout) AddForce(f Force) {
	l.forces = append(l.forces, f)
}

// Nodes returns nodes information.
func (l *Layout) Positions() map[string]*Object {
	return l.objects
}

// Links returns graph data links.
func (l *Layout) Links() []*graph.Link {
	return l.links
}
