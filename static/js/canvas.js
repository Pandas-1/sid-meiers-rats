const GRID_SIZE = 30
const CELL_SIZE = 16 // 800px / 50 cells
const imageCache = {}
const TILE_W = 38.4
const TILE_H = 21.6
const ORIGIN_X = GRID_SIZE * TILE_W / 2
const ORIGIN_Y = 40
const grassA = new Image()
grassA.src = '/static/images/terrain/grass_a.png'
const grassB = new Image()
grassB.src = '/static/images/terrain/grass_b.png'


function drawTerrainTile(ctx, gx, gy) {
    const p = toScreen(gx, gy)
    const img = (gx + gy) % 2 === 0 ? grassA : grassB
    if (!img.complete) return
    ctx.drawImage(img, p.x - TILE_W / 2, p.y, TILE_W, TILE_H)
}


function toScreen(gridX, gridY) {
    return {
        x: (gridX - gridY) * (TILE_W / 2) + ORIGIN_X,
        y: (gridX + gridY) * (TILE_H / 2) + ORIGIN_Y,
    }
}

function toGrid(screenX, screenY) {
    const dx = screenX - ORIGIN_X
    const dy = screenY - ORIGIN_Y
    return {
        x: Math.floor((dx / (TILE_W / 2) + dy / (TILE_H / 2)) / 2),
        y: Math.floor((dy / (TILE_H / 2) - dx / (TILE_W / 2)) / 2),
    }
}

function initCanvas() {
    const canvas = document.getElementById('gameCanvas')
    canvas.width = GRID_SIZE * TILE_W
    canvas.height = GRID_SIZE * TILE_H + ORIGIN_Y * 2
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
    for (let gx = 0; gx < GRID_SIZE; gx++) {
        for (let gy = 0; gy < GRID_SIZE; gy++) {
            drawTerrainTile(ctx, gx, gy)
        }
    }
    ctx.strokeStyle = '#2a2a4a'
    ctx.lineWidth = 0.5
    for (let x = 0; x <= GRID_SIZE; x++) {
        const a = toScreen(x, 0)
        const b = toScreen(x, GRID_SIZE)
        ctx.beginPath()
        ctx.moveTo(a.x, a.y)
        ctx.lineTo(b.x, b.y)
        ctx.stroke()
    }
    for (let y = 0; y <= GRID_SIZE; y++) {
        const a = toScreen(0, y)
        const b = toScreen(GRID_SIZE, y)
        ctx.beginPath()
        ctx.moveTo(a.x, a.y)
        ctx.lineTo(b.x, b.y)
        ctx.stroke()
    }
}

function drawBuilding(ctx, building) {
    const north = toScreen(building.GridX, building.GridY)
    const w = (building.Width + building.Height) * (TILE_W / 2)
    const h = (building.Width + building.Height) * (TILE_H / 2)

    const drawX = north.x - w / 2
    const drawY = north.y

    const img = getImage(building.Name)
    if (img.complete) {
        ctx.drawImage(img, drawX, drawY, w, h)
    } else {
        ctx.fillStyle = '#32475e'
        ctx.fillRect(drawX, drawY, w, h)
        img.onload = () => ctx.drawImage(img, drawX, drawY, w, h)
    }

    ctx.fillStyle = '#ffffff'
    ctx.font = '8px sans-serif'
    ctx.fillText(building.Name, drawX + 2, drawY + 10)
}