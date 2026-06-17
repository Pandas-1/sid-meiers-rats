const GRID_SIZE = 50
const CELL_SIZE = 16 // 800px / 50 cells

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

function drawBuilding(ctx, building) {
    const x = building.GridX * CELL_SIZE
    const y = building.GridY * CELL_SIZE
    const w = building.Width * CELL_SIZE
    const h = building.Height * CELL_SIZE

    ctx.fillStyle = '#4a90d9'
    ctx.fillRect(x, y, w, h)
    ctx.strokeStyle = '#ffffff'
    ctx.strokeRect(x, y, w, h)
}