package main

import "fmt"

// PieceType represents the different bugs in Hive.
type PieceType string

const (
	Queen       PieceType = "Queen"
	Ant         PieceType = "Ant"
	Beetle      PieceType = "Beetle"
	Grasshopper PieceType = "Grasshopper"
	Spider      PieceType = "Spider"
)

// Color represents the player owning a piece.
type Color string

const (
	White Color = "White"
	Black Color = "Black"
)

// Piece represents an individual bug on the board or in a hand.
type Piece struct {
	Type  PieceType `json:"type"`
	Color Color     `json:"color"`
}

// GameState holds the master state of the game.
type GameState struct {
	// Grid maps hex coordinate strings ("q,r") to a stack of pieces.
	// Using string keys ensures stable JSON serialization.
	Grid map[string][]Piece `json:"grid"`

	// Hands tracks the pieces each player has yet to place.
	Hands map[Color][]Piece `json:"hands"`

	// CurrentTurn indicates whose turn it is (White or Black).
	CurrentTurn Color `json:"current_turn"`

	// GameStatus could be "active", "white_win", "black_win", or "draw".
	GameStatus string `json:"game_status"`

	// TurnCount tracks the number of turns each player has taken.
	TurnCount map[Color]int `json:"turn_count"`
}

// NewGame initializes a fresh game with all pieces in each player's hand.
func NewGame() *GameState {
	return &GameState{
		Grid: make(map[string][]Piece),
		Hands: map[Color][]Piece{
			White: initialHand(White),
			Black: initialHand(Black),
		},
		CurrentTurn: White,
		GameStatus:  "active",
		TurnCount: map[Color]int{
			White: 0,
			Black: 0,
		},
	}
}

// initialHand returns the standard set of 11 Hive pieces for a player.
func initialHand(color Color) []Piece {
	return []Piece{
		{Type: Queen, Color: color},
		{Type: Beetle, Color: color}, {Type: Beetle, Color: color},
		{Type: Grasshopper, Color: color}, {Type: Grasshopper, Color: color}, {Type: Grasshopper, Color: color},
		{Type: Spider, Color: color}, {Type: Spider, Color: color},
		{Type: Ant, Color: color}, {Type: Ant, Color: color}, {Type: Ant, Color: color},
	}
}

// GetTopPiece returns the piece on top of a stack at a given hex.
func (gs *GameState) GetTopPiece(h Hex) *Piece {
	key := fmt.Sprintf("%d,%d", h.Q, h.R)
	stack, exists := gs.Grid[key]
	if !exists || len(stack) == 0 {
		return nil
	}
	return &stack[len(stack)-1]
}

// PlacePiece adds a piece to the grid at the specified hex.
func (gs *GameState) PlacePiece(h Hex, p Piece) {
	key := fmt.Sprintf("%d,%d", h.Q, h.R)
	gs.Grid[key] = append(gs.Grid[key], p)
}

// IsValidPlacement checks if a piece can be placed at a hex.
func (gs *GameState) IsValidPlacement(color Color, pt PieceType, h Hex) bool {
	// 1. Correct turn?
	if color != gs.CurrentTurn {
		return false
	}

	// 2. Already occupied?
	if gs.GetTopPiece(h) != nil {
		return false
	}

	// 3. Piece in hand?
	hand := gs.Hands[color]
	hasPiece := false
	for _, p := range hand {
		if p.Type == pt {
			hasPiece = true
			break
		}
	}
	if !hasPiece {
		return false
	}

	// 4. Queen must be placed by 4th turn.
	if gs.TurnCount[color] == 3 {
		hasQueen := false
		for _, p := range hand {
			if p.Type == Queen {
				hasQueen = true
				break
			}
		}
		if hasQueen && pt != Queen {
			return false
		}
	}

	// 5. Adjacency rules
	gridSize := len(gs.Grid)

	// First piece: any hex is valid.
	if gridSize == 0 {
		return true
	}

	neighbors := h.Neighbors()
	hasFriendlyNeighbor := false
	hasEnemyNeighbor := false

	for _, n := range neighbors {
		neighborPiece := gs.GetTopPiece(n)
		if neighborPiece != nil {
			if neighborPiece.Color == color {
				hasFriendlyNeighbor = true
			} else {
				hasEnemyNeighbor = true
			}
		}
	}

	// Second piece: must be adjacent to the first piece (regardless of color).
	if gridSize == 1 {
		return hasFriendlyNeighbor || hasEnemyNeighbor
	}

	// Subsequent pieces: must be adjacent to friendly AND NOT adjacent to enemy.
	return hasFriendlyNeighbor && !hasEnemyNeighbor
}

// PlayPiece places a piece from a player's hand onto the grid and switches the turn.
// Returns false if the piece is not in the player's hand or the move is invalid.
func (gs *GameState) PlayPiece(h Hex, pt PieceType) bool {
	if !gs.IsValidPlacement(gs.CurrentTurn, pt, h) {
		return false
	}

	hand := gs.Hands[gs.CurrentTurn]
	foundIdx := -1
	for i, p := range hand {
		if p.Type == pt {
			foundIdx = i
			break
		}
	}

	// (We already checked in IsValidPlacement, but for safety)
	if foundIdx == -1 {
		return false
	}

	// Remove from hand
	p := hand[foundIdx]
	gs.Hands[gs.CurrentTurn] = append(hand[:foundIdx], hand[foundIdx+1:]...)

	// Add to grid
	gs.PlacePiece(h, p)

	// Increment turn count
	gs.TurnCount[gs.CurrentTurn]++

	// Toggle turn
	if gs.CurrentTurn == White {
		gs.CurrentTurn = Black
	} else {
		gs.CurrentTurn = White
	}

	return true
}
