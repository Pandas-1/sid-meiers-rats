const GRID_SIZE = 50
const CELL_SIZE = 16
const imageCache = {}
const TILE_W = 24
const TILE_H = 12
const ORIGIN_X = GRID_SIZE * TILE_W / 2
const ORIGIN_Y = 40

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

function drawGrid(ctx) {
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

function drawHPBar(ctx, x, y, w, currentHP, maxHP) {
    const barH = 3
    const pct = currentHP / maxHP
    ctx.fillStyle = '#ff0000'
    ctx.fillRect(x, y - barH - 1, w, barH)
    ctx.fillStyle = '#00ff00'
    ctx.fillRect(x, y - barH - 1, w * pct, barH)
}

function drawBuilding(ctx, building) {
    const north = toScreen(building.x, building.y)
    const w = (building.width + building.height) * (TILE_W / 2)
    const h = (building.width + building.height) * (TILE_H / 2)

    const drawX = north.x - w / 2
    const drawY = north.y

    const img = getImage(building.name)
    if (img.complete) {
        ctx.drawImage(img, drawX, drawY, w, h)
    } else {
        ctx.fillStyle = '#32475e'
        ctx.fillRect(drawX, drawY, w, h)
        img.onload = () => ctx.drawImage(img, drawX, drawY, w, h)
    }

    ctx.fillStyle = '#ffffff'
    ctx.font = '8px sans-serif'
    ctx.fillText(building.name, drawX + 2, drawY + 10)

    drawHPBar(ctx, drawX, drawY, w, building.current_hp, building.max_hp)
}

function drawTroop(ctx, t) {
    const { x, y } = toScreen(t.x, t.y)
    const size = TILE_W * 0.5

    ctx.fillStyle = '#ff6b6b'
    ctx.beginPath()
    ctx.arc(x, y - size / 2, size / 2, 0, Math.PI * 2)
    ctx.fill()

    ctx.fillStyle = '#ffffff'
    ctx.font = '6px sans-serif'
    ctx.fillText(t.name[0], x - 2, y - size / 2 + 2)

    drawHPBar(ctx, x - size / 2, y - size, size, t.current_hp, t.max_hp)
}

function renderBattle(ctx, state) {
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    drawGrid(ctx)
    if (state.buildings) state.buildings.forEach(b => drawBuilding(ctx, b))
    if (state.troops) state.troops.forEach(t => drawTroop(ctx, t))
}
function getImage(name) {
    if (imageCache[name]) return imageCache[name]
    const img = new Image()
    img.src = `/static/images/buildings/${name.toLowerCase().replace(/ /g, '_')}.png`
    imageCache[name] = img
    console.log("image called for?")
    return img
}