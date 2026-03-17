# Project: bug-of-war (2D Turn-Based Strategy Game)

## Overview
We are building a 2D, flat-UI, web-based multiplayer board game inspired by "Hive". The game has no board; pieces are placed adjacently to form a continuous grid. 
I will be hosting this locally on my machine, acting as both the server and a client, with my friend connecting via a secure tunnel (e.g., Ngrok).

## Tech Stack
* **Backend:** Go (using `gorilla/websocket` for real-time bidirectional communication).
* **Frontend:** HTML5 Canvas API (for rendering the flat 2D hexagonal tiles) and vanilla JavaScript.
* **State Management:** In-memory on the Go server. No database is required for this MVP.

## Current Project State
- **Hex Grid System:** Fully implemented with axial/cube coordinate support and mouse-to-hex rounding.
- **Game State:** Implemented with support for piece types (Queen, Ant, Beetle, Grasshopper, Spider) and piece stacking (Beetle).
- **Frontend Rendering:** Basic HTML5 Canvas rendering is in place, fetching state from the server via a JSON API.
- **Testing Infrastructure:** Unit tests exist for hex math and the core game state.

## Roadmap
1. [x] Hex grid math and basic canvas rendering.
2. [x] Core game state and piece definitions.
3. [x] Placement Interaction & Rules (Selecting from hand, placement validation, turn switching).
4. [ ] **Next:** Movement Interaction & Rules (Moving on-board pieces, "One Hive" rule, specific bug patterns).
5. [ ] WebSocket infrastructure for real-time multiplayer.
6. [ ] Win condition detection (Queen surrounded).

## Architecture & Data Structures
* **Grid System:** The game uses a flat-topped hexagonal grid. Use cube coordinates ($x, y, z$ where $x + y + z = 0$) or axial coordinates ($q, r$) for all piece placements and distance calculations.
* **Game State:** The server holds the master state: current turn, piece locations, unplayed pieces in hand, and game status (active/win/draw).
* **Communication:** * Client sends: \`PLACE_PIECE\`, \`MOVE_PIECE\`.
    * Server broadcasts: \`STATE_UPDATE\`, \`ERROR\` (for invalid moves).

## Core Game Rules (The Mechanics)
1.  **The One Hive Rule:** The pieces in play must ALWAYS form a single, contiguous cluster. A piece cannot move if its movement would temporarily or permanently split the hive.
2.  **Placement:** New pieces must touch your own color and cannot touch the opponent's color (except for the very first placement of the game).
3.  **The Queen:** Must be placed within the first four turns. No pieces can move until the Queen is placed.
4.  **Winning:** The game ends when a Queen is completely surrounded by 6 pieces (of any color).

## Piece Movement Rules
* **Queen:** Moves exactly 1 space per turn.
* **Beetle:** Moves 1 space per turn, but can climb on top of the hive (stacking on other pieces).
* **Grasshopper:** Jumps over a straight line of contiguous pieces to the next empty space.
* **Spider:** Moves exactly 3 spaces crawling around the outside perimeter of the hive. Cannot backtrack.
* **Ant:** Can crawl anywhere around the outside perimeter of the hive.

## Agent Directives
* Write clean, modular code. Separate the WebSocket networking logic from the game validation logic.
* Do not generate massive monolithic files. Break things out (e.g., \`hex_math.go\`, \`rules.go\`, \`server.go\`).
* Prioritize robust backend validation. The client UI should be "dumb" and only render what the server dictates.

## Development Workflow
* **Mandatory Verification:** Before completing any code changes or sub-tasks, Gemini MUST run \`go test -v .\` to confirm no regressions were introduced.
* **Context Preservation:** Gemini should check for feature-specific context in \`.gemini/<branch-name>/context.md\` if it exists.
