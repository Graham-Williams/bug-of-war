# BUG-OF-WAR

A 2D turn-based strategy game inspired by "Hive".

## Current Progress: Hex Grid Math & Rendering
The project now includes:
*   **Hexagonal Math**: Core logic for axial/cube coordinates, distance, neighbors, and rounding in `hex_math.go`.
*   **Interactive Canvas**: A flat-topped hexagonal grid rendered via HTML5 Canvas in `app.js`.

## Prerequisites

* **Go Version 1.26.1**

## Getting Started

1. **Install Dependencies:**
   Run the following command to download and tidy up project dependencies:
   ```bash
   go mod tidy
   ```

2. **Run the Server:**
   ```bash
   go run .
   ```

3. **Access the Game:**
   Open your browser and navigate to:
   [http://localhost:8080](http://localhost:8080)

## TODO
[ ] WebSocket infrastructure
[ ] Initial game state & piece types
[ ] Movement validation logic
[ ] Add license
