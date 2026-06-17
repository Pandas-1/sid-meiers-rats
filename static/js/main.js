const ctx = initCanvas()
console.log("Canvas initialized")
drawGrid(ctx)

const token = localStorage.getItem('token')
if (!token) {
    window.location.href = '/static/login.html'
}

function render(buildings) {
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    drawGrid(ctx)
    buildings.forEach(b => drawBuilding(ctx, b))
}

async function loadVillage() {
    const response = await fetch('/village', {
        headers: { 'Authorization': token }
    })
    if (!response.ok) {
        console.log("Auth failed:", response.status)
        return
    }
    const buildings = await response.json()
    console.log("buildings", buildings)
    render(buildings)
}
loadVillage()