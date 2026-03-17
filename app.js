const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');
const coordsDisplay = document.getElementById('coords');
const turnDisplay = document.getElementById('turn-display');
const statusDisplay = document.getElementById('status-display');
const whiteHandDiv = document.getElementById('white-hand');
const blackHandDiv = document.getElementById('black-hand');

const HEX_SIZE = 40; // Distance from center to vertex

// --- Colors & Aesthetics ---
const COLORS = {
    White: { fill: '#ffffff', text: '#000000' },
    Black: { fill: '#333333', text: '#ffffff' },
    Empty: { fill: '#eee4d3', stroke: '#cbbba0' },
    Selected: { stroke: '#ffcc00', width: 4 },
    Hover: { stroke: '#a8d5ff', width: 2 }
};

const PIECE_LABELS = {
    Queen: 'Q',
    Ant: 'A',
    Beetle: 'B',
    Grasshopper: 'G',
    Spider: 'S'
};

// --- State ---
let gameState = { grid: {}, hands: {}, current_turn: 'White', game_status: 'active' };
let hoveredHex = { q: null, r: null };
let selectedFromHand = null; // { type, color }
let selectedFromBoard = null; // { q, r }

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
function drawHex(q, r, stack = [], isHovered = false, isSelected = false) {
    const center = axialToPixel(q, r);
    const piece = stack.length > 0 ? stack[stack.length - 1] : null;
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
    if (isSelected) {
        ctx.strokeStyle = COLORS.Selected.stroke;
        ctx.lineWidth = COLORS.Selected.width;
    } else if (isHovered) {
        ctx.strokeStyle = COLORS.Hover.stroke;
        ctx.lineWidth = COLORS.Hover.width;
    } else {
        ctx.strokeStyle = COLORS.Empty.stroke;
        ctx.lineWidth = 1;
    }
    ctx.stroke();

    // Draw piece label if exists
    if (piece) {
        ctx.fillStyle = colorTheme.text;
        ctx.font = 'bold 20px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(PIECE_LABELS[piece.type] || '?', center.x, center.y);

        // Stack indicator
        if (stack.length > 1) {
            ctx.fillStyle = colorTheme.text;
            ctx.font = 'bold 12px sans-serif';
            ctx.fillText(stack.length, center.x + 15, center.y - 15);
            
            // Draw a second ring for stack
            ctx.beginPath();
            ctx.arc(center.x, center.y, HEX_SIZE * 0.8, 0, Math.PI * 2);
            ctx.strokeStyle = colorTheme.text;
            ctx.lineWidth = 1;
            ctx.stroke();
        }
    }

    // Small coordinate debug text
    ctx.fillStyle = piece ? colorTheme.text : '#aaa';
    ctx.font = '9px sans-serif';
    ctx.fillText(`${q},${r}`, center.x, center.y + (piece ? 15 : 4));
}

function updateHandUI() {
    ['White', 'Black'].forEach(color => {
        const div = color === 'White' ? whiteHandDiv : blackHandDiv;
        const hand = gameState.hands[color] || [];
        
        const counts = hand.reduce((acc, p) => {
            acc[p.type] = (acc[p.type] || 0) + 1;
            return acc;
        }, {});

        div.innerHTML = '';
        Object.keys(PIECE_LABELS).forEach(type => {
            const count = counts[type] || 0;
            const btn = document.createElement('div');
            const isSelected = selectedFromHand?.type === type && selectedFromHand?.color === color;
            btn.className = `piece-button ${count === 0 ? 'disabled' : ''} ${isSelected ? 'selected' : ''}`;
            btn.innerText = PIECE_LABELS[type];
            btn.title = `${type} (${count})`;
            
            if (count > 0 && color === gameState.current_turn && gameState.game_status === 'active') {
                btn.onclick = () => {
                    selectedFromBoard = null;
                    if (isSelected) {
                        selectedFromHand = null;
                    } else {
                        selectedFromHand = { type, color };
                    }
                    updateHandUI();
                    draw();
                };
            }
            
            const badge = document.createElement('span');
            badge.style.position = 'absolute';
            badge.style.fontSize = '8px';
            badge.style.bottom = '-5px';
            badge.style.right = '-5px';
            badge.style.background = '#007bff';
            badge.style.color = 'white';
            badge.style.borderRadius = '50%';
            badge.style.width = '12px';
            badge.style.height = '12px';
            badge.style.display = count > 1 ? 'flex' : 'none';
            badge.style.justifyContent = 'center';
            badge.style.alignItems = 'center';
            badge.innerText = count;
            btn.style.position = 'relative';
            btn.appendChild(badge);

            div.appendChild(btn);
        });
    });
}

function draw() {
    if (!canvas.width || !canvas.height) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Dynamic range based on pieces
    let minQ = -3, maxQ = 3, minR = -3, maxR = 3;
    Object.keys(gameState.grid).forEach(key => {
        const [q, r] = key.split(',').map(Number);
        minQ = Math.min(minQ, q - 2);
        maxQ = Math.max(maxQ, q + 2);
        minR = Math.min(minR, r - 2);
        maxR = Math.max(maxR, r + 2);
    });

    for (let q = minQ; q <= maxQ; q++) {
        for (let r = minR; r <= maxR; r++) {
            const key = `${q},${r}`;
            const isHovered = q === hoveredHex.q && r === hoveredHex.r;
            const isSelected = selectedFromBoard && q === selectedFromBoard.q && r === selectedFromBoard.r;
            const stack = gameState.grid[key] || [];
            drawHex(q, r, stack, isHovered, isSelected);
        }
    }
}

// --- Data Fetching ---
async function fetchState() {
    try {
        const response = await fetch('/state');
        const newState = await response.json();
        
        // Only update if something changed (to avoid flickering if possible)
        gameState = newState;
        updateUI();
        draw();
    } catch (err) {
        console.error("Failed to fetch game state:", err);
    }
}

function updateUI() {
    turnDisplay.innerText = `Turn: ${gameState.current_turn}`;
    const statusMap = {
        'active': 'Game In Progress',
        'white_win': 'White Wins!',
        'black_win': 'Black Wins!',
        'draw': 'Draw!'
    };
    statusDisplay.innerText = statusMap[gameState.game_status] || gameState.game_status;
    statusDisplay.style.color = gameState.game_status === 'active' ? '#666' : '#dc3545';
    updateHandUI();
}

async function placePiece(q, r, type) {
    try {
        const response = await fetch('/place', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ q, r, type })
        });
        
        if (response.ok) {
            gameState = await response.json();
            selectedFromHand = null;
            updateUI();
            draw();
        } else {
            const err = await response.text();
            alert("Invalid move: " + err);
        }
    } catch (err) {
        console.error("Failed to place piece:", err);
    }
}

async function movePiece(fromQ, fromR, toQ, toR) {
    try {
        const response = await fetch('/move', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ from_q: fromQ, from_r: fromR, to_q: toQ, to_r: toR })
        });
        
        if (response.ok) {
            gameState = await response.json();
            selectedFromBoard = null;
            updateUI();
            draw();
        } else {
            const err = await response.text();
            alert("Invalid move: " + err);
        }
    } catch (err) {
        console.error("Failed to move piece:", err);
    }
}

// --- Setup & Interaction ---
function resize() {
    canvas.width = window.innerWidth * 0.95;
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
    coordsDisplay.innerText = `Q: ${hex.q}, R: ${hex.r}`;
    
    if (hex.q !== hoveredHex.q || hex.r !== hoveredHex.r) {
        hoveredHex = hex;
        draw();
    }
});

canvas.addEventListener('click', (e) => {
    if (gameState.game_status !== 'active') return;

    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const hex = pixelToAxial(x, y);
    const key = `${hex.q},${hex.r}`;
    const stack = gameState.grid[key] || [];
    const topPiece = stack.length > 0 ? stack[stack.length - 1] : null;

    if (selectedFromHand) {
        placePiece(hex.q, hex.r, selectedFromHand.type);
    } else if (selectedFromBoard) {
        if (selectedFromBoard.q === hex.q && selectedFromBoard.r === hex.r) {
            // Deselect
            selectedFromBoard = null;
        } else {
            // Attempt move
            movePiece(selectedFromBoard.q, selectedFromBoard.r, hex.q, hex.r);
        }
    } else if (topPiece && topPiece.color === gameState.current_turn) {
        // Select from board
        selectedFromBoard = { q: hex.q, r: hex.r };
        selectedFromHand = null;
    }
    
    updateHandUI();
    draw();
});

// Initial fetch
fetchState();
setInterval(fetchState, 5000); 
