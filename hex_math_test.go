package main

import "testing"

func TestHexToCube(t *testing.T) {
	h := Hex{Q: 1, R: 2}
	c := h.ToCube()
	if c.X != 1 || c.Y != -3 || c.Z != 2 {
		t.Errorf("Expected Cube{1, -3, 2}, got %+v", c)
	}
}

func TestCubeToHex(t *testing.T) {
	c := Cube{X: 1, Y: -3, Z: 2}
	h := c.ToHex()
	if h.Q != 1 || h.R != 2 {
		t.Errorf("Expected Hex{1, 2}, got %+v", h)
	}
}

func TestDistance(t *testing.T) {
	h1 := Hex{Q: 0, R: 0}
	h2 := Hex{Q: 2, R: 0}
	if d := h1.Distance(h2); d != 2 {
		t.Errorf("Expected distance 2, got %d", d)
	}

	h3 := Hex{Q: -1, R: -1}
	if d := h1.Distance(h3); d != 2 {
		t.Errorf("Expected distance 2, got %d", d)
	}
}

func TestNeighbors(t *testing.T) {
	h := Hex{Q: 0, R: 0}
	n := h.Neighbors()
	if len(n) != 6 {
		t.Errorf("Expected 6 neighbors, got %d", len(n))
	}

	// Neighbor at {1, 0}
	found := false
	for _, neighbor := range n {
		if neighbor.Q == 1 && neighbor.R == 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected neighbor {1, 0} not found")
	}
}

func TestHexRounding(t *testing.T) {
	// Near (0,0)
	h := HexRounding(0.1, 0.1)
	if h.Q != 0 || h.R != 0 {
		t.Errorf("Expected Hex{0, 0}, got %+v", h)
	}

	// Near (1, -2)
	h = HexRounding(0.9, -1.9)
	if h.Q != 1 || h.R != -2 {
		t.Errorf("Expected Hex{1, -2}, got %+v", h)
	}
}
