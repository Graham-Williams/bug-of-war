package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var game = NewGame()

func main() {
	// API to get the current game state
	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	})

	// API to place a piece from hand
	http.HandleFunc("/place", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Q    int       `json:"q"`
			R    int       `json:"r"`
			Type PieceType `json:"type"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		h := Hex{Q: req.Q, R: req.R}
		if !game.PlayPiece(h, req.Type) {
			http.Error(w, "Invalid move", http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	})

	// API to move a piece on the board
	http.HandleFunc("/move", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			FromQ int `json:"from_q"`
			FromR int `json:"from_r"`
			ToQ   int `json:"to_q"`
			ToR   int `json:"to_r"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		from := Hex{Q: req.FromQ, R: req.FromR}
		to := Hex{Q: req.ToQ, R: req.ToR}
		if !game.MovePiece(from, to) {
			http.Error(w, "Invalid move", http.StatusUnprocessableEntity)
			return
		}

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
