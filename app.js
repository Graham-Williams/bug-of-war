const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');
const coordsDisplay = document.getElementById('coords');

const HEX_SIZE = 40; // Distance from center to vertex

// --- Colors & Aesthetics ---
const COLORS = {
    White: { fill: '#ffffff', text: '#000000' },
    Black: { fill: '#333333', text: '#ffffff' },
    Empty: { fill: '#eee4d3', stroke: '#cbbba0' },
    Selected: { stroke: '#a8d5ff', width: 4 }
};

const PIECE_LABELS = {
    Queen: 'Q',
    Ant: 'A',
    Beetle: 'B',
    Grasshopper: 'G',
    Spider: 'S'
};

// --- State ---
let gameState = { grid: {}, hands: {}, current_turn: 'White' };
let selectedHex = { q: null, r: null };

// --- Hex Math Helpers ---
function getHexCorner(center, size, i) {
    const angleDeg = 60 * i;
    const angleRad = (Math.PI / 180) * angleDeg;
    return {
        x: center.x + size * Math.cos(angleRad),
        y: center.y + size * Math.sin(angleRad)
    };
}

function axialToPixel(q, r) {
    const x = HEX_SIZE * (3/2 * q);
    const y = HEX_SIZE * (Math.sqrt(3)/2 * q + Math.sqrt(3) * r);
    return { x: x + canvas.width / 2, y: y + canvas.height / 2 };
}

function pixelToAxial(px, py) {
    const x = px - canvas.width / 2;
    const y = py - canvas.height / 2;
    const q = (2/3 * x) / HEX_SIZE;
    const r = (-1/3 * x + Math.sqrt(3)/3 * y) / HEX_SIZE;
    return hexRound(q, r);
}

function hexRound(q, r) {
    let x = q, z = r, y = -x - z;
    let rx = Math.round(x), ry = Math.round(y), rz = Math.round(z);
    const xDiff = Math.abs(rx - x), yDiff = Math.abs(ry - y), zDiff = Math.abs(rz - z);
    if (xDiff > yDiff && xDiff > zDiff) rx = -ry - rz;
    else if (yDiff > zDiff) ry = -rx - rz;
    else rz = -rx - ry;
    return { q: rx, r: rz };
}

// --- Rendering ---
function drawHex(q, r, piece = null, isSelected = false) {
    const center = axialToPixel(q, r);
    const colorTheme = piece ? COLORS[piece.color] : COLORS.Empty;

    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
        const corner = getHexCorner(center, HEX_SIZE, i);
        if (i === 0) ctx.moveTo(corner.x, corner.y);
        else ctx.lineTo(corner.x, corner.y);
    }
    ctx.closePath();

    // Fill hex
    ctx.fillStyle = colorTheme.fill;
    ctx.fill();

    // Stroke hex
    ctx.strokeStyle = isSelected ? COLORS.Selected.stroke : (COLORS.Empty.stroke);
    ctx.lineWidth = isSelected ? COLORS.Selected.width : 2;
    ctx.stroke();

    // Draw piece label if exists
    if (piece) {
        ctx.fillStyle = colorTheme.text;
        ctx.font = 'bold 20px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(PIECE_LABELS[piece.type] || '?', center.x, center.y);
    }

    // Small coordinate debug text
    ctx.fillStyle = piece ? colorTheme.text : '#aaa';
    ctx.font = '9px sans-serif';
    ctx.fillText(`${q},${r}`, center.x, center.y + (piece ? 15 : 4));
}

function draw() {
    if (!canvas.width || !canvas.height) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // 1. Draw background grid
    const range = 5;
    for (let q = -range; q <= range; q++) {
        for (let r = -range; r <= range; r++) {
            if (Math.abs(q + r) <= range) {
                const key = `${q},${r}`;
                // Only draw empty hex if no piece is there
                if (!gameState.grid[key]) {
                    const isSelected = q === selectedHex.q && r === selectedHex.r;
                    drawHex(q, r, null, isSelected);
                }
            }
        }
    }

    // 2. Draw pieces from the server state
    for (const [coords, stack] of Object.entries(gameState.grid)) {
        const [q, r] = coords.split(',').map(Number);
        const topPiece = stack[stack.length - 1];
        const isSelected = q === selectedHex.q && r === selectedHex.r;
        drawHex(q, r, topPiece, isSelected);
    }
}

// --- Data Fetching ---
async function fetchState() {
    try {
        const response = await fetch('/state');
        const newState = await response.json();
        
        // Convert the backend map key {Q:0, R:0} to a string key "0,0" for easier lookup
        // The backend JSON for a map with struct keys is usually an object with stringified keys like "{\"Q\":0,\"R\":0}"
        // BUT Go's JSON encoder for maps with struct keys is only supported if they implement TextMarshaler.
        // Let's assume for now we might need to adjust the backend if the format is tricky.
        
        // Actual fix for Go map JSON: it's better to use string keys on the backend or 
        // handle the complex key parsing here. For simplicity in this step, let's update 
        // the backend to use string keys in a future turn if needed.
        
        gameState = newState;
        draw();
    } catch (err) {
        console.error("Failed to fetch game state:", err);
    }
}

// --- Setup & Interaction ---
function resize() {
    canvas.width = window.innerWidth * 0.9;
    canvas.height = window.innerHeight * 0.8;
    draw();
}

window.addEventListener('resize', resize);
resize();

canvas.addEventListener('mousemove', (e) => {
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const hex = pixelToAxial(x, y);
    coordsDisplay.innerText = `Q: ${hex.q}, R: ${hex.r} | Turn: ${gameState.current_turn}`;
    
    if (hex.q !== selectedHex.q || hex.r !== selectedHex.r) {
        selectedHex = hex;
        draw();
    }
});

// Initial fetch
fetchState();
setInterval(fetchState, 2000); // Poll every 2s for now
