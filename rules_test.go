package main

import (
	"testing"
)

func TestIsValidPlacement(t *testing.T) {
	gs := NewGame()

	// 1. First move (White) - Any hex should be valid (usually 0,0)
	h0 := Hex{0, 0}
	if !gs.IsValidPlacement(White, Ant, h0) {
		t.Errorf("First move (White) at (0,0) should be valid")
	}
	gs.PlayPiece(h0, Ant) // White plays Ant at (0,0)

	// 2. Second move (Black) - Must be adjacent to the first piece
	h1 := Hex{1, 0}
	if !gs.IsValidPlacement(Black, Queen, h1) {
		t.Errorf("Second move (Black) adjacent to first piece should be valid")
	}
	
	hFar := Hex{5, 5}
	if gs.IsValidPlacement(Black, Queen, hFar) {
		t.Errorf("Second move (Black) NOT adjacent to first piece should be invalid")
	}
	gs.PlayPiece(h1, Queen) // Black plays Queen at (1,0)

	// 3. Third move (White) - Must be adjacent to friendly, NOT adjacent to enemy
	h2 := Hex{-1, 0} // Adjacent to White Ant (0,0), Not adjacent to Black Queen (1,0)
	if !gs.IsValidPlacement(White, Queen, h2) {
		t.Errorf("White placement adjacent to friendly and NOT enemy should be valid")
	}

	hInvalid := Hex{1, -1} // Adjacent to White Ant (0,0) AND Black Queen (1,0)
	if gs.IsValidPlacement(White, Queen, hInvalid) {
		t.Errorf("White placement adjacent to enemy piece should be invalid")
	}
	gs.PlayPiece(h2, Queen) // White plays Queen at (-1,0)

	// 4. Occupied Hex - Should be invalid
	if gs.IsValidPlacement(Black, Ant, h0) {
		t.Errorf("Placement on occupied hex (0,0) should be invalid")
	}
}

func TestQueenPlacementRule(t *testing.T) {
	gs := NewGame()
	
	// White plays 3 pieces (not Queen)
	if !gs.PlayPiece(Hex{0, 0}, Ant) { t.Fatalf("White 1 failed") }
	if !gs.PlayPiece(Hex{1, 0}, Ant) { t.Fatalf("Black 1 failed") }
	if !gs.PlayPiece(Hex{-1, 0}, Ant) { t.Fatalf("White 2 failed") }
	if !gs.PlayPiece(Hex{2, 0}, Ant) { t.Fatalf("Black 2 failed") }
	if !gs.PlayPiece(Hex{-1, 1}, Ant) { t.Fatalf("White 3 failed") }
	if !gs.PlayPiece(Hex{3, 0}, Ant) { t.Fatalf("Black 3 failed") }

	// Now it's White's 4th turn. White MUST play Queen.
	if gs.IsValidPlacement(White, Ant, Hex{-2, 1}) {
		t.Errorf("White must place Queen by 4th turn; playing Ant should be invalid")
	}
	if !gs.IsValidPlacement(White, Queen, Hex{-2, 1}) {
		t.Errorf("White must place Queen by 4th turn; playing Queen should be valid")
	}
}
