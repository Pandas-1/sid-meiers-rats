const GRID_SIZE = 50
const CELL_SIZE = 16 // 800px / 50 cells
const imageCache = {}

function initCanvas() {
    const canvas = document.getElementById('gameCanvas')
    canvas.width = GRID_SIZE * CELL_SIZE
    canvas.height = GRID_SIZE * CELL_SIZE
    return canvas.getContext('2d')
}

function getImage(name) {
    if (imageCache[name]) return imageCache[name]
    const img = new Image()
    img.src = `/static/images/buildings/${name.toLowerCase().replace(/ /g, '_')}.png`
    imageCache[name] = img
    return img
}

function drawGrid(ctx) {
    ctx.strokeStyle = '#2a2a4a'
    ctx.lineWidth = 0.5
    for (let x = 0; x <= GRID_SIZE; x++) {
        ctx.beginPath()
        ctx.moveTo(x * CELL_SIZE, 0)
        ctx.lineTo(x * CELL_SIZE, GRID_SIZE * CELL_SIZE)
        ctx.stroke()
    }
    for (let y = 0; y <= GRID_SIZE; y++) {
        ctx.beginPath()
        ctx.moveTo(0, y * CELL_SIZE)
        ctx.lineTo(GRID_SIZE * CELL_SIZE, y * CELL_SIZE)
        ctx.stroke()
    }
}

function drawBuilding(ctx, building) {
    const x = building.GridX * CELL_SIZE
    const y = building.GridY * CELL_SIZE
    const w = building.Width * CELL_SIZE
    const h = building.Height * CELL_SIZE

    const img = getImage(building.Name)
    if (img.complete) {
        ctx.drawImage(img, x, y, w, h)
    } else {
        // fallback color while image loads
        ctx.fillStyle = '#32475e'
        ctx.fillRect(x, y, w, h)
        img.onload = () => ctx.drawImage(img, x, y, w, h)
    }

    // building name label
    ctx.fillStyle = '#ffffff'
    ctx.font = '8px sans-serif'
    ctx.fillText(building.Name, x + 2, y + 10)
}