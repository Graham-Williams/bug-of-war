package main

import "math"

// Hex represents a single hexagonal tile using axial coordinates (q, r).
// In a flat-topped grid, q is the column and r is the row.
type Hex struct {
	Q, R int
}

// Cube represents hexagonal coordinates in cube space (x, y, z).
// x + y + z must always equal 0.
type Cube struct {
	X, Y, Z int
}

// ToCube converts axial coordinates (q, r) to cube coordinates (x, y, z).
func (h Hex) ToCube() Cube {
	return Cube{
		X: h.Q,
		Y: -h.Q - h.R,
		Z: h.R,
	}
}

// ToHex converts cube coordinates (x, y, z) to axial coordinates (q, r).
func (c Cube) ToHex() Hex {
	return Hex{
		Q: c.X,
		R: c.Z,
	}
}

// Distance calculates the Manhattan distance between two hexes in the grid.
func (h Hex) Distance(other Hex) int {
	a := h.ToCube()
	b := other.ToCube()
	return (abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)) / 2
}

// Neighbors returns the 6 immediate neighbors of a hex.
func (h Hex) Neighbors() []Hex {
	directions := []Hex{
		{1, 0}, {1, -1}, {0, -1},
		{-1, 0}, {-1, 1}, {0, 1},
	}
	neighbors := make([]Hex, 6)
	for i, d := range directions {
		neighbors[i] = Hex{h.Q + d.Q, h.R + d.R}
	}
	return neighbors
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// HexRounding handles fractional coordinates to find the nearest integer hex.
// This is useful for converting mouse clicks (pixel -> axial) to the grid.
func HexRounding(fq, fr float64) Hex {
	fx := fq
	fz := fr
	fy := -fx - fz

	rx := math.Round(fx)
	ry := math.Round(fy)
	rz := math.Round(fz)

	xDiff := math.Abs(rx - fx)
	yDiff := math.Abs(ry - fy)
	zDiff := math.Abs(rz - fz)

	if xDiff > yDiff && xDiff > zDiff {
		rx = -ry - rz
	} else if yDiff > zDiff {
		ry = -rx - rz
	} else {
		rz = -rx - ry
	}

	return Hex{int(rx), int(rz)}
}
