const GRID_SIZE = 50
const CELL_SIZE = 16

function initCanvas() {
    const canvas = document.getElementById('gameCanvas')
    canvas.width = GRID_SIZE * CELL_SIZE
    canvas.height = GRID_SIZE * CELL_SIZE
    return canvas.getContext('2d')
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

function drawHPBar(ctx, x, y, w, currentHP, maxHP) {
    const barH = 3
    const pct = currentHP / maxHP
    ctx.fillStyle = '#ff0000'
    ctx.fillRect(x, y - barH - 1, w, barH)
    ctx.fillStyle = '#00ff00'
    ctx.fillRect(x, y - barH - 1, w * pct, barH)
}

function drawBuilding(ctx, b) {
    const x = b.x * CELL_SIZE
    const y = b.y * CELL_SIZE
    const w = b.width * CELL_SIZE
    const h = b.height * CELL_SIZE

    ctx.fillStyle = '#4a90d9'
    ctx.fillRect(x, y, w, h)
    ctx.strokeStyle = '#ffffff'
    ctx.strokeRect(x, y, w, h)

    ctx.fillStyle = '#ffffff'
    ctx.font = '7px sans-serif'
    ctx.fillText(b.name, x + 2, y + 10)

    drawHPBar(ctx, x, y, w, b.current_hp, b.max_hp)
}

function drawTroop(ctx, t) {
    const x = t.x * CELL_SIZE
    const y = t.y * CELL_SIZE
    const size = CELL_SIZE * 0.8

    ctx.fillStyle = '#ff6b6b'
    ctx.beginPath()
    ctx.arc(x, y, size / 2, 0, Math.PI * 2)
    ctx.fill()

    ctx.fillStyle = '#ffffff'
    ctx.font = '6px sans-serif'
    ctx.fillText(t.name[0], x - 2, y + 2)

    drawHPBar(ctx, x - size/2, y - size/2, size, t.current_hp, t.max_hp)
}

function renderBattle(ctx, state) {
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    drawGrid(ctx)
    if (state.buildings) state.buildings.forEach(b => drawBuilding(ctx, b))
    if (state.troops) state.troops.forEach(t => drawTroop(ctx, t))
}