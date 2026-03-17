package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var game = NewGame()

func init() {
	// TEMPORARY: Pre-place a variety of pieces for visual verification.
	// This block will be removed once the MOVE_PIECE/PLACE_PIECE logic is implemented.
	game.PlacePiece(Hex{Q: 0, R: 0}, Piece{Type: Queen, Color: White})
	game.PlacePiece(Hex{Q: 1, R: 0}, Piece{Type: Queen, Color: Black})
	game.PlacePiece(Hex{Q: 0, R: 1}, Piece{Type: Ant, Color: White})
	game.PlacePiece(Hex{Q: -1, R: 1}, Piece{Type: Beetle, Color: Black})
	game.PlacePiece(Hex{Q: 1, R: -1}, Piece{Type: Grasshopper, Color: White})
	game.PlacePiece(Hex{Q: 2, R: -1}, Piece{Type: Spider, Color: Black})
}

func main() {
	// API to get the current game state
	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	})

	// Serve static files from the current directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
