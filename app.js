const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');
const coordsDisplay = document.getElementById('coords');

const HEX_SIZE = 40; // Distance from center to vertex

// --- State ---
let selectedHex = { q: null, r: null };

// --- Hex Math Helpers ---
function getHexCorner(center, size, i) {
    const angleDeg = 60 * i; // Flat-topped hex starts at 0 degrees
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
    let x = q;
    let z = r;
    let y = -x - z;

    let rx = Math.round(x);
    let ry = Math.round(y);
    let rz = Math.round(z);

    const xDiff = Math.abs(rx - x);
    const yDiff = Math.abs(ry - y);
    const zDiff = Math.abs(rz - z);

    if (xDiff > yDiff && xDiff > zDiff) {
        rx = -ry - rz;
    } else if (yDiff > zDiff) {
        ry = -rx - rz;
    } else {
        rz = -rx - ry;
    }

    return { q: rx, r: rz };
}

// --- Rendering ---
function drawHex(q, r, color = '#eee', stroke = '#ccc') {
    const center = axialToPixel(q, r);
    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
        const corner = getHexCorner(center, HEX_SIZE, i);
        if (i === 0) ctx.moveTo(corner.x, corner.y);
        else ctx.lineTo(corner.x, corner.y);
    }
    ctx.closePath();
    ctx.fillStyle = color;
    ctx.fill();
    ctx.strokeStyle = stroke;
    ctx.lineWidth = 2;
    ctx.stroke();

    // Draw coordinates for debugging
    ctx.fillStyle = '#aaa';
    ctx.font = '10px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText(`${q},${r}`, center.x, center.y + 4);
}

function draw() {
    if (!canvas.width || !canvas.height) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Draw a small 5x5 grid for testing
    const range = 5;
    for (let q = -range; q <= range; q++) {
        for (let r = -range; r <= range; r++) {
            if (Math.abs(q + r) <= range) {
                const isSelected = q === selectedHex.q && r === selectedHex.r;
                drawHex(q, r, isSelected ? '#a8d5ff' : '#f9f9f9');
            }
        }
    }
}

// --- Setup & Sizing ---
function resize() {
    canvas.width = window.innerWidth * 0.9;
    canvas.height = window.innerHeight * 0.8;
    console.log(`Canvas resized to: ${canvas.width}x${canvas.height}`);
    draw();
}

window.addEventListener('resize', resize);
resize();

// --- Interaction ---
canvas.addEventListener('mousedown', (e) => {
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    const hex = pixelToAxial(x, y);
    console.log(`Hex clicked: Q: ${hex.q}, R: ${hex.r}`);
});

canvas.addEventListener('mousemove', (e) => {
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    const hex = pixelToAxial(x, y);
    coordsDisplay.innerText = `Q: ${hex.q}, R: ${hex.r}`;
    
    if (hex.q !== selectedHex.q || hex.r !== selectedHex.r) {
        selectedHex = hex;
        draw();
    }
});
