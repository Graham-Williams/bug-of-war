# BUG-OF-WAR

A 2D turn-based strategy game inspired by "Hive".

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

## Testing

To run the full suite of Go unit tests:
```bash
go test -v .
```

## Roadmap
[x] Hex Grid & Canvas Rendering
[x] Initial game state & piece types
[x] Placement Interaction & Rules
[ ] Movement validation logic
[ ] WebSocket infrastructure
[ ] Win condition detection
[ ] Add license
