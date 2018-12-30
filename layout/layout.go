package layout

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/divan/graphx/graph"
)

// stableThreshold determines the movement diff needed to
// call the system stable
const stableThreshold = 2.001

// Layout implements Layout interface for force-directed 3D graph.
type Layout struct {
	g *graph.Graph

	objects map[string]*Object // node ID as a key
	keys    []string           // IDs in the order of adding
	links   []*graph.Link

	confMu sync.RWMutex
	config Config
	forces []Force
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// New creates a new layout from the given config.
func New(g *graph.Graph, config Config) *Layout {
	forces := forcesFromConfig(config)
	l := NewWithForces(g, forces...)
	l.config = config
	return l
}

// NewWithForces initializes layout with data and custom set of forces.
func NewWithForces(g *graph.Graph, forces ...Force) *Layout {
	l := &Layout{
		g:       g,
		objects: make(map[string]*Object),
		keys:    make([]string, 0, len(g.Nodes())),
		links:   g.Links(),
		forces:  forces,
	}

	l.initPositions()

	return l
}

// initPositions inits layout graph from the original graph data.
func (l *Layout) initPositions() {
	for _, node := range l.g.Nodes() {
		l.addObject(node)
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
	lastIdx := len(l.keys) // use last item index for calculating pseudo-random positions

	// TODO: handle weight

	x, y, z := randomPosition(lastIdx)
	object := NewObjectID(x, y, z, node.ID())

	l.objects[node.ID()] = object
	l.keys = append(l.keys, node.ID())
}

// randomPosition generates x,y,z coordinates pseudo-randomly spread around the
// spherical surface.
func randomPosition(i int) (x, y, z float64) {
	radius := 10 * math.Cbrt(float64(i))
	rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
	yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

	x = radius * math.Cos(rollAngle)
	y = radius * math.Sin(rollAngle)
	z = radius * math.Sin(yawAngle)
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
	for i := 0; i < n; i++ {
		l.UpdatePositions()
	}
}

// UpdatePositions recalculates nodes' positions, applying all the forces.
// It returns average amount of movement generated by this step.
func (l *Layout) UpdatePositions() float64 {
	l.resetForces()

	l.confMu.RLock()
	forces := l.forces
	l.confMu.RUnlock()
	for _, force := range forces {
		apply := force.Rule()
		apply(force, l.objects, l.links)
	}

	return l.integrate()
}

func (l *Layout) resetForces() {
	for i := range l.objects {
		l.objects[i].force = ZeroForce()
	}
}

// AddForce adds force to the internal list of forces.
// FIXME: this breaks sync between config and forces
func (l *Layout) AddForce(f Force) {
	l.confMu.Lock()
	defer l.confMu.Unlock()
	l.forces = append(l.forces, f)
}

// Positions returns nodes information as a map with ID key.
// TODO(divan): rename it to something else after fixing legacy usage.
func (l *Layout) Positions() map[string]*Object {
	return l.objects
}

// Positions returns nodes information as a slice, where index order is equal to the
// original graph nodes order.
func (l *Layout) PositionsSlice() []*Position {
	ret := make([]*Position, len(l.keys))
	for i, id := range l.keys {
		obj := l.objects[id]
		pos := &Position{obj._X, obj._Y, obj._Z}
		ret[i] = pos
	}
	return ret
}

// Links returns graph data links.
func (l *Layout) Links() []*graph.Link {
	return l.links
}

// Config returns current config.
func (l *Layout) Config() Config {
	l.confMu.RLock()
	defer l.confMu.RUnlock()
	return l.config
}

// SetConfig updates current config and forces.
func (l *Layout) SetConfig(c Config) {
	l.confMu.Lock()
	l.config = c
	l.forces = forcesFromConfig(c)
	l.confMu.Unlock()
}

// SetPositions overwrites objects positions and recalculates layout internal stuff
// to be in sync with new positions.
// Positions slice should be the same size and order as Nodes.
func (l *Layout) SetPositions(positions []*Position) {
	// recalculate objects with new positions
	//l.resetObjects()
	for i, node := range l.g.Nodes() {
		id := node.ID()
		pos := positions[i]
		obj := l.objects[id]
		obj.SetPosition(pos.X, pos.Y, pos.Z)
	}
}

func (l *Layout) resetObjects() {
	l.objects = make(map[string]*Object)
}

// Graph returns original data graph for layout.
func (l *Layout) Graph() *graph.Graph {
	return l.g
}
