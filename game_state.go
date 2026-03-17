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

// GetOccupiedHexes returns a list of all hexes that have at least one piece.
func (gs *GameState) GetOccupiedHexes() []Hex {
	hexes := make([]Hex, 0, len(gs.Grid))
	for key := range gs.Grid {
		var q, r int
		fmt.Sscanf(key, "%d,%d", &q, &r)
		hexes = append(hexes, Hex{Q: q, R: r})
	}
	return hexes
}

// IsHiveContiguous checks if all pieces on the board form a single connected group.
func (gs *GameState) IsHiveContiguous(hypotheticalRemoved *Hex) bool {
	occupied := gs.GetOccupiedHexes()
	
	remaining := make(map[string]bool)
	var startKey string

	for _, h := range occupied {
		key := fmt.Sprintf("%d,%d", h.Q, h.R)
		if hypotheticalRemoved != nil && h.Q == hypotheticalRemoved.Q && h.R == hypotheticalRemoved.R {
			if len(gs.Grid[key]) > 1 {
				remaining[key] = true
				if startKey == "" {
					startKey = key
				}
				continue
			}
			continue
		}
		
		remaining[key] = true
		if startKey == "" {
			startKey = key
		}
	}

	if len(remaining) <= 1 {
		return true
	}

	// BFS to count reachable hexes
	visited := make(map[string]bool)
	queue := []string{startKey}
	visited[startKey] = true
	count := 0

	for len(queue) > 0 {
		currKey := queue[0]
		queue = queue[1:]
		count++

		var cq, cr int
		fmt.Sscanf(currKey, "%d,%d", &cq, &cr)
		current := Hex{Q: cq, R: cr}

		for _, neighbor := range current.Neighbors() {
			nKey := fmt.Sprintf("%d,%d", neighbor.Q, neighbor.R)
			if remaining[nKey] && !visited[nKey] {
				visited[nKey] = true
				queue = append(queue, nKey)
			}
		}
	}

	return count == len(remaining)
}

// HasPlacedQueen checks if a player has their Queen on the board.
func (gs *GameState) HasPlacedQueen(color Color) bool {
	for _, stack := range gs.Grid {
		for _, p := range stack {
			if p.Color == color && p.Type == Queen {
				return true
			}
		}
	}
	return false
}

// CheckWinCondition checks for surrounded Queens and updates GameStatus.
func (gs *GameState) CheckWinCondition() {
	var whiteSurrounded, blackSurrounded bool

	for key, stack := range gs.Grid {
		if len(stack) == 0 {
			continue
		}
		top := stack[len(stack)-1]
		if top.Type == Queen {
			var q, r int
			fmt.Sscanf(key, "%d,%d", &q, &r)
			h := Hex{Q: q, R: r}
			
			neighbors := h.Neighbors()
			count := 0
			for _, n := range neighbors {
				if len(gs.Grid[fmt.Sprintf("%d,%d", n.Q, n.R)]) > 0 {
					count++
				}
			}

			if count == 6 {
				if top.Color == White {
					whiteSurrounded = true
				} else {
					blackSurrounded = true
				}
			}
		}
	}

	if whiteSurrounded && blackSurrounded {
		gs.GameStatus = "draw"
	} else if whiteSurrounded {
		gs.GameStatus = "black_win"
	} else if blackSurrounded {
		gs.GameStatus = "white_win"
	}
}

// IsValidMove checks if moving a piece from 'from' to 'to' is legal.
func (gs *GameState) IsValidMove(color Color, from, to Hex) bool {
	if color != gs.CurrentTurn {
		return false
	}

	// 1. Queen must be placed to move any piece.
	if !gs.HasPlacedQueen(color) {
		return false
	}

	// 2. Must have a piece at 'from'.
	piece := gs.GetTopPiece(from)
	if piece == nil || piece.Color != color {
		return false
	}

	// 3. 'to' must be different from 'from'.
	if from.Q == to.Q && from.R == to.R {
		return false
	}

	// 4. One Hive Rule: hive must be contiguous during and after move.
	if !gs.IsHiveContiguous(&from) {
		return false
	}

	// Destination must be adjacent to the remaining hive.
	isAdjacentToHive := false
	for _, n := range to.Neighbors() {
		// Don't count the 'from' hex if it will be empty
		if n.Q == from.Q && n.R == from.R {
			if len(gs.Grid[fmt.Sprintf("%d,%d", from.Q, from.R)]) > 1 {
				isAdjacentToHive = true
				break
			}
			continue
		}
		if len(gs.Grid[fmt.Sprintf("%d,%d", n.Q, n.R)]) > 0 {
			isAdjacentToHive = true
			break
		}
	}
	if !isAdjacentToHive {
		return false
	}

	// 5. Common Rule: Pieces cannot move into occupied spaces (except Beetles).
	targetStack := gs.Grid[fmt.Sprintf("%d,%d", to.Q, to.R)]
	if len(targetStack) > 0 && piece.Type != Beetle {
		return false
	}

	// 6. Piece-specific rules
	switch piece.Type {
	case Queen:
		// Queen moves 1 step.
		isNeighbor := false
		for _, n := range from.Neighbors() {
			if n.Q == to.Q && n.R == to.R {
				isNeighbor = true
				break
			}
		}
		if !isNeighbor {
			return false
		}
		// Sliding rule (basic check)
		commonNeighbors := 0
		for _, n1 := range from.Neighbors() {
			for _, n2 := range to.Neighbors() {
				if n1.Q == n2.Q && n1.R == n2.R {
					if len(gs.Grid[fmt.Sprintf("%d,%d", n1.Q, n1.R)]) > 0 {
						commonNeighbors++
					}
				}
			}
		}
		if commonNeighbors == 2 {
			return false
		}

	case Beetle:
		// Beetle moves 1 step. Can climb on others.
		isNeighbor := false
		for _, n := range from.Neighbors() {
			if n.Q == to.Q && n.R == to.R {
				isNeighbor = true
				break
			}
		}
		if !isNeighbor {
			return false
		}
		// Beetle can also climb, so sliding rule is slightly different 
		// (can only climb if it's not "trapped" by even taller stacks, 
		// but for now 1 step is enough).

	case Ant:
		// Ant moves any distance around the perimeter.
		if len(targetStack) > 0 {
			return false
		}

		// Use BFS to find all reachable perimeter hexes
		visited := make(map[string]bool)
		queue := []Hex{from}
		visited[fmt.Sprintf("%d,%d", from.Q, from.R)] = true
		
		reachable := false
		toKey := fmt.Sprintf("%d,%d", to.Q, to.R)

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			currKey := fmt.Sprintf("%d,%d", curr.Q, curr.R)
			if currKey == toKey && (curr.Q != from.Q || curr.R != from.R) {
				reachable = true
				break
			}

			for _, neighbor := range curr.Neighbors() {
				nKey := fmt.Sprintf("%d,%d", neighbor.Q, neighbor.R)
				if visited[nKey] {
					continue
				}

				// 1. Cannot be occupied
				if len(gs.Grid[nKey]) > 0 {
					continue
				}

				// 2. Must be adjacent to the hive (excluding the Ant itself)
				isAdjacentToHive := false
				for _, nn := range neighbor.Neighbors() {
					if nn.Q == from.Q && nn.R == from.R {
						if len(gs.Grid[fmt.Sprintf("%d,%d", from.Q, from.R)]) > 1 {
							isAdjacentToHive = true
							break
						}
						continue
					}
					if len(gs.Grid[fmt.Sprintf("%d,%d", nn.Q, nn.R)]) > 0 {
						isAdjacentToHive = true
						break
					}
				}

				if !isAdjacentToHive {
					continue
				}

				// 3. Sliding rule
				occupiedCommon := 0
				for _, cn1 := range curr.Neighbors() {
					for _, cn2 := range neighbor.Neighbors() {
						if cn1.Q == cn2.Q && cn1.R == cn2.R {
							if len(gs.Grid[fmt.Sprintf("%d,%d", cn1.Q, cn1.R)]) > 0 {
								// Count as occupied if it's in Grid AND it's NOT our current ant 'from' hex
								if cn1.Q != from.Q || cn1.R != from.R || len(gs.Grid[fmt.Sprintf("%d,%d", from.Q, from.R)]) > 1 {
									occupiedCommon++
								}
							}
						}
					}
				}
				if occupiedCommon == 2 {
					continue
				}

				visited[nKey] = true
				queue = append(queue, neighbor)
			}
		}

		return reachable


	case Grasshopper:
		// Grasshopper jumps in a straight line over 1 or more pieces to the first empty space.
		if len(targetStack) > 0 {
			return false
		}

		// 1. Must be in a straight line
		diffQ := to.Q - from.Q
		diffR := to.R - from.R

		// In a hex grid, a straight line means one coordinate is constant, 
		// or the sum is constant (cube: x+y+z=0).
		// Directions: {1, 0}, {1, -1}, {0, -1}, {-1, 0}, {-1, 1}, {0, 1}
		dirQ, dirR := 0, 0
		if diffQ == 0 {
			dirR = 1
			if diffR < 0 { dirR = -1 }
		} else if diffR == 0 {
			dirQ = 1
			if diffQ < 0 { dirQ = -1 }
		} else if diffQ == -diffR {
			dirQ = 1
			if diffQ < 0 { dirQ = -1 }
			dirR = -dirQ
		} else {
			// Not a straight line
			return false
		}

		// 2. Must jump over at least one piece
		steps := 0
		currQ, currR := from.Q+dirQ, from.R+dirR
		for currQ != to.Q || currR != to.R {
			if len(gs.Grid[fmt.Sprintf("%d,%d", currQ, currR)]) == 0 {
				// Gap found in the jump - invalid
				return false
			}
			steps++
			currQ += dirQ
			currR += dirR
		}

		if steps == 0 {
			// Didn't jump over anything
			return false
		}
		return true

	case Spider:
		// Spider moves exactly 3 steps around the perimeter of the hive.
		if len(targetStack) > 0 {
			return false
		}

		fromKey := fmt.Sprintf("%d,%d", from.Q, from.R)

		// Use BFS to find all possible destinations after exactly 3 steps
		type spiderPath struct {
			hex  Hex
			path []string // to avoid backtracking
		}

		queue := []spiderPath{{hex: from, path: []string{fromKey}}}
		reachedInLess := make(map[string]bool)
		reachedInLess[fromKey] = true
		validDestinations := make(map[string]bool)

		for steps := 0; steps < 3; steps++ {
			nextQueue := []spiderPath{}
			for _, p := range queue {
				for _, neighbor := range p.hex.Neighbors() {
					nKey := fmt.Sprintf("%d,%d", neighbor.Q, neighbor.R)
					
					if len(gs.Grid[nKey]) > 0 { continue }

					backtracked := false
					for _, prev := range p.path {
						if prev == nKey {
							backtracked = true
							break
						}
					}
					if backtracked { continue }

					isAdjacentToHive := false
					for _, nn := range neighbor.Neighbors() {
						if nn.Q == from.Q && nn.R == from.R {
							if len(gs.Grid[fmt.Sprintf("%d,%d", from.Q, from.R)]) > 1 {
								isAdjacentToHive = true
								break
							}
							continue
						}
						if len(gs.Grid[fmt.Sprintf("%d,%d", nn.Q, nn.R)]) > 0 {
							isAdjacentToHive = true
							break
						}
					}
					if !isAdjacentToHive { continue }

					occupiedCommon := 0
					for _, cn1 := range p.hex.Neighbors() {
						for _, cn2 := range neighbor.Neighbors() {
							if cn1.Q == cn2.Q && cn1.R == cn2.R {
								if len(gs.Grid[fmt.Sprintf("%d,%d", cn1.Q, cn1.R)]) > 0 {
									if cn1.Q != from.Q || cn1.R != from.R || len(gs.Grid[fmt.Sprintf("%d,%d", from.Q, from.R)]) > 1 {
										occupiedCommon++
									}
								}
							}
						}
					}
					if occupiedCommon == 2 { continue }

					newPath := append([]string{}, p.path...)
					newPath = append(newPath, nKey)
					nextQueue = append(nextQueue, spiderPath{hex: neighbor, path: newPath})
					
					if steps < 2 {
						reachedInLess[nKey] = true
					} else {
						validDestinations[nKey] = true
					}
				}
			}
			queue = nextQueue
		}

		// Remove destinations that were reachable in fewer steps
		for k := range reachedInLess {
			delete(validDestinations, k)
		}

		toKey := fmt.Sprintf("%d,%d", to.Q, to.R)
		if !validDestinations[toKey] {
			return false
		}
		return true

	}

	return true
}

// MovePiece moves the top piece from 'from' to 'to'.
func (gs *GameState) MovePiece(from, to Hex) bool {
	if !gs.IsValidMove(gs.CurrentTurn, from, to) {
		return false
	}

	fromKey := fmt.Sprintf("%d,%d", from.Q, from.R)
	toKey := fmt.Sprintf("%d,%d", to.Q, to.R)

	stack := gs.Grid[fromKey]
	piece := stack[len(stack)-1]
	
	// Remove from 'from'
	gs.Grid[fromKey] = stack[:len(stack)-1]
	if len(gs.Grid[fromKey]) == 0 {
		delete(gs.Grid, fromKey)
	}

	// Add to 'to'
	gs.Grid[toKey] = append(gs.Grid[toKey], piece)

	// Check win conditions
	gs.CheckWinCondition()

	// Toggle turn
	if gs.CurrentTurn == White {
		gs.CurrentTurn = Black
	} else {
		gs.CurrentTurn = White
	}

	return true
}

// PlayPiece places a piece from a player's hand onto the grid and switches the turn.
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

	if foundIdx == -1 {
		return false
	}

	p := hand[foundIdx]
	gs.Hands[gs.CurrentTurn] = append(hand[:foundIdx], hand[foundIdx+1:]...)

	gs.PlacePiece(h, p)

	// Increment turn count
	gs.TurnCount[gs.CurrentTurn]++

	// Check win conditions
	gs.CheckWinCondition()

	// Toggle turn
	if gs.CurrentTurn == White {
		gs.CurrentTurn = Black
	} else {
		gs.CurrentTurn = White
	}

	return true
}
