package main

import (
	"fmt"
	"testing"
)

func TestNewGame(t *testing.T) {
	gs := NewGame()

	if len(gs.Hands[White]) != 11 {
		t.Errorf("Expected White to have 11 pieces, got %d", len(gs.Hands[White]))
	}

	if len(gs.Hands[Black]) != 11 {
		t.Errorf("Expected Black to have 11 pieces, got %d", len(gs.Hands[Black]))
	}

	if gs.CurrentTurn != White {
		t.Errorf("Expected current turn to be White, got %s", gs.CurrentTurn)
	}

	if len(gs.Grid) != 0 {
		t.Errorf("Expected empty grid, got %d pieces", len(gs.Grid))
	}
}

func TestPlacePiece(t *testing.T) {
	gs := NewGame()
	h := Hex{Q: 0, R: 0}
	p := Piece{Type: Queen, Color: White}

	gs.PlacePiece(h, p)

	key := fmt.Sprintf("%d,%d", h.Q, h.R)
	if len(gs.Grid[key]) != 1 {
		t.Errorf("Expected 1 piece at (0,0), got %d", len(gs.Grid[key]))
	}

	top := gs.GetTopPiece(h)
	if top == nil || top.Type != Queen || top.Color != White {
		t.Errorf("Expected Queen on top, got %v", top)
	}
}

func TestStacking(t *testing.T) {
	gs := NewGame()
	h := Hex{Q: 1, R: 1}
	p1 := Piece{Type: Ant, Color: White}
	p2 := Piece{Type: Beetle, Color: Black}

	gs.PlacePiece(h, p1)
	gs.PlacePiece(h, p2)

	key := fmt.Sprintf("%d,%d", h.Q, h.R)
	if len(gs.Grid[key]) != 2 {
		t.Errorf("Expected 2 pieces at (1,1), got %d", len(gs.Grid[key]))
	}

	top := gs.GetTopPiece(h)
	if top == nil || top.Type != Beetle || top.Color != Black {
		t.Errorf("Expected Beetle on top, got %v", top)
	}
}

func TestPlayPiece(t *testing.T) {
	gs := NewGame()
	h := Hex{Q: 0, R: 0}

	// 1. White plays a Queen
	success := gs.PlayPiece(h, Queen)
	if !success {
		t.Errorf("Expected PlayPiece to succeed")
	}

	// Verify grid
	top := gs.GetTopPiece(h)
	if top == nil || top.Type != Queen || top.Color != White {
		t.Errorf("Expected Queen at (0,0), got %v", top)
	}

	// Verify hand removal
	if len(gs.Hands[White]) != 10 {
		t.Errorf("Expected White to have 10 pieces left, got %d", len(gs.Hands[White]))
	}

	// Verify turn toggled
	if gs.CurrentTurn != Black {
		t.Errorf("Expected current turn to be Black, got %s", gs.CurrentTurn)
	}

	// 2. Black plays a Spider
	h2 := Hex{Q: 1, R: 0}
	success = gs.PlayPiece(h2, Spider)
	if !success {
		t.Errorf("Expected Black to play Spider successfully")
	}

	// Verify grid
	top = gs.GetTopPiece(h2)
	if top == nil || top.Type != Spider || top.Color != Black {
		t.Errorf("Expected Spider at (1,0), got %v", top)
	}

	// Verify hand removal
	if len(gs.Hands[Black]) != 10 {
		t.Errorf("Expected Black to have 10 pieces left, got %d", len(gs.Hands[Black]))
	}

	// Verify turn toggled
	if gs.CurrentTurn != White {
		t.Errorf("Expected current turn to be White, got %s", gs.CurrentTurn)
	}
}
