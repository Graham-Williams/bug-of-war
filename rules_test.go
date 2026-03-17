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

func TestQueenMovement(t *testing.T) {
	gs := NewGame()

	// 1. Setup board with Queens
	gs.PlayPiece(Hex{0, 0}, Queen) // White Queen
	gs.PlayPiece(Hex{1, 0}, Queen) // Black Queen

	// 2. White Queen tries to move (should be legal)
	if !gs.IsValidMove(White, Hex{0, 0}, Hex{0, 1}) {
		t.Errorf("White Queen should be allowed to move 1 step to (0,1)")
	}
	if !gs.MovePiece(Hex{0, 0}, Hex{0, 1}) {
		t.Fatalf("White Queen move failed")
	}

	// 3. Black Queen tries to move (should be legal)
	if !gs.IsValidMove(Black, Hex{1, 0}, Hex{1, 1}) {
		t.Errorf("Black Queen should be allowed to move 1 step to (1,1)")
	}
	if !gs.MovePiece(Hex{1, 0}, Hex{1, 1}) {
		t.Fatalf("Black Queen move failed")
	}

	// 4. White Queen tries to move 2 steps (should be illegal)
	if gs.IsValidMove(White, Hex{0, 1}, Hex{0, 3}) {
		t.Errorf("White Queen moving 2 steps should be invalid")
	}
}

func TestOneHiveRule(t *testing.T) {
	gs := NewGame()
	
	// Setup: (-1,0)W Ant - (0,0)W Queen - (1,0)B Queen
	if !gs.PlayPiece(Hex{0, 0}, Queen) { t.Fatalf("White Queen failed") }
	if !gs.PlayPiece(Hex{1, 0}, Queen) { t.Fatalf("Black Queen failed") }
	if !gs.PlayPiece(Hex{-1, 0}, Ant) { t.Fatalf("White Ant failed") }

	// Set turn back to White
	gs.CurrentTurn = White
	// Moving the middle piece (0,0) would split the hive
	if gs.IsValidMove(White, Hex{0, 0}, Hex{0, 1}) {
		t.Errorf("Moving White Queen at (0,0) should be invalid as it splits the hive")
	}

	// Moving the end piece (-1,0) should be valid
	if !gs.IsValidMove(White, Hex{-1, 0}, Hex{-1, 1}) {
		t.Errorf("Moving White Ant at (-1,0) should be valid as it doesn't split the hive")
	}
}

func TestWinCondition(t *testing.T) {
	gs := NewGame()
	
	// Surround White Queen at (0,0)
	gs.PlayPiece(Hex{0, 0}, Queen) // White (0,0)
	gs.PlayPiece(Hex{1, 0}, Queen) // Black (1,0)
	
	// Black places 5 more pieces around (0,0)
	// (Note: we use PlayPiece but skip IsValidPlacement checks for speed here)
	gs.Grid["1,-1"] = []Piece{{Type: Ant, Color: Black}}
	gs.Grid["0,-1"] = []Piece{{Type: Ant, Color: Black}}
	gs.Grid["-1,0"] = []Piece{{Type: Ant, Color: Black}}
	gs.Grid["-1,1"] = []Piece{{Type: Ant, Color: Black}}
	gs.Grid["0,1"] = []Piece{{Type: Ant, Color: Black}}
	
	gs.CheckWinCondition()
	if gs.GameStatus != "black_win" {
		t.Errorf("GameStatus should be black_win after White Queen is surrounded, got %s", gs.GameStatus)
	}
}

func TestGrasshopperMovement(t *testing.T) {
	gs := NewGame()

	// 1. Setup board with Queen and Grasshopper
	gs.PlayPiece(Hex{0, 0}, Queen)      // White Queen
	gs.PlayPiece(Hex{1, 0}, Queen)      // Black Queen
	gs.PlayPiece(Hex{-1, 0}, Grasshopper) // White GH

	// 2. White GH tries to move to (0,1) - not a straight line
	if gs.IsValidMove(White, Hex{-1, 0}, Hex{0, 1}) {
		t.Errorf("GH move should be illegal when not in a straight line")
	}

	// 3. White GH tries to jump over White Queen and Black Queen to (2,0)
	gs.CurrentTurn = White
	if !gs.IsValidMove(White, Hex{-1, 0}, Hex{2, 0}) {
		t.Errorf("White GH should be allowed to jump over 2 pieces to (2,0)")
	}

	// 4. White GH tries to move over 3 pieces
	gs.CurrentTurn = Black
	if !gs.PlayPiece(Hex{2, 0}, Ant) { t.Fatalf("Black Ant failed") }
	
	// Now try to jump from -1 over (0, 1, 2) to 3
	gs.CurrentTurn = White
	if !gs.IsValidMove(White, Hex{-1, 0}, Hex{3, 0}) {
		t.Errorf("White GH should be allowed to jump over 3 pieces to (3,0)")
	}

	// Now introduce a gap at (1,0)
	delete(gs.Grid, "1,0")
	if gs.IsValidMove(White, Hex{-1, 0}, Hex{3, 0}) {
		t.Errorf("GH should NOT be allowed to jump over a gap")
	}
}

func TestSpiderMovement(t *testing.T) {
	gs := NewGame()

	// 1. Setup board
	gs.PlayPiece(Hex{0, 0}, Queen) // White
	gs.PlayPiece(Hex{1, 0}, Queen) // Black
	gs.PlayPiece(Hex{-1, 0}, Spider) // White

	// 2. White Spider tries to move
	gs.CurrentTurn = White
	// Possible destinations from (-1,0):
	// Step 1: (-1,1), (0,-1)
	// Step 2 from (-1,1): (0,1)
	// Step 3 from (0,1): (1,1)
	
	if !gs.IsValidMove(White, Hex{-1, 0}, Hex{1, 1}) {
		t.Errorf("Spider should be allowed to move exactly 3 steps to (1,1)")
	}

	if gs.IsValidMove(White, Hex{-1, 0}, Hex{0, 1}) {
		t.Errorf("Spider should NOT be allowed to move only 2 steps")
	}
}

func TestAntMovement(t *testing.T) {
	gs := NewGame()

	// 1. Setup board
	gs.PlayPiece(Hex{0, 0}, Queen) // White
	gs.PlayPiece(Hex{1, 0}, Queen) // Black
	gs.PlayPiece(Hex{-1, 0}, Ant)   // White

	// 2. White Ant tries to move
	gs.CurrentTurn = White
	// Should be able to reach any perimeter hex: (1,1), (0,1), (-1,1), (0,-1), (1,-1), (2,0), (2,-1)
	if !gs.IsValidMove(White, Hex{-1, 0}, Hex{2, 0}) {
		t.Errorf("Ant should be allowed to move any distance around perimeter to (2,0)")
	}

	if gs.IsValidMove(White, Hex{-1, 0}, Hex{5, 5}) {
		t.Errorf("Ant should NOT be allowed to move to a non-perimeter hex")
	}
}
