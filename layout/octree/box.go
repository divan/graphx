package octree

import (
	"math"
)

// Box represents bounding box for octant.
type Box struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
	Front  float64
	Back   float64
}

// NewBox creates new bounding box.
func NewBox(left, right, top, bottom, front, back float64) *Box {
	return &Box{
		Left:   left,
		Right:  right,
		Bottom: bottom,
		Top:    top,
		Front:  front,
		Back:   back,
	}
}

// NewBoxPoints creates new bounding box from two points.
func NewBoxPoints(from, to Point) *Box {
	dx := math.Abs(from.X() - to.X())
	dy := math.Abs(from.Y() - to.Y())
	dz := math.Abs(from.Z() - to.Z())
	maxSide := math.Max(dx, math.Max(dy, dz))
	box := &Box{}
	box.Left = math.Min(from.X(), to.X())
	box.Right = box.Left + maxSide
	box.Bottom = math.Min(from.Y(), to.Y())
	box.Top = box.Bottom + maxSide
	box.Front = math.Min(from.Z(), to.Z())
	box.Back = box.Front + maxSide
	return box
}

// NewZeroBox returns new bounding box with zero as a first point.
func NewZeroBox(p Point) *Box {
	zero := NewPoint("", 0, 0, 0, 0)
	return NewBoxPoints(zero, p)
}

// Center returns the middle point of the box.
func (b *Box) Center() (float64, float64, float64) {
	x := b.Left + b.Width()/2
	y := b.Bottom + b.Height()/2
	z := b.Front + b.Depth()/2
	return x, y, z
}

func (b *Box) Width() float64  { return b.Right - b.Left }
func (b *Box) Height() float64 { return b.Top - b.Bottom }
func (b *Box) Depth() float64  { return b.Back - b.Front }
