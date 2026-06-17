const ctx = initCanvas()
console.log("Canvas initialized")
drawGrid(ctx)

const token = localStorage.getItem('token')
if (!token) {
    window.location.href = '/static/login.html'
}

let mode = null
let selectedInstanceID = null
let buildings = []
const buildingSelect = document.getElementById('buildingSelect')
const status = document.getElementById('status')

async function loadShop() {
    const res = await fetch('/shop/buildings', {
        headers: {'Authorization': token}
    })
    const shopBuildings = await res.json()
    shopBuildings.forEach( b => {
        if (b.Name === "Town Hall"){
            return;
        }
        const option = document.createElement("option")
        option.value = b.BuildingID
        option.textContent = b.Name + " " +b.CostResource1 + "gold " + b.CostResource2 + "elixer"
        buildingSelect.appendChild(option)
        console.log(b)

    })
}

ctx.canvas.addEventListener('click', async function(e) {
    const rect = ctx.canvas.getBoundingClientRect()
    const x = Math.floor((e.clientX - rect.left) / CELL_SIZE)
    const y = Math.floor((e.clientY - rect.top) / CELL_SIZE)

    if (mode === 'place') {
        const buildingID = parseInt(buildingSelect.value)
        const res = await fetch('/village/place', {
            method: 'POST',
            headers: {
                'Authorization': token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ building_id: buildingID, X: x, Y: y })
        })
        if (res.ok) {
            status.textContent = 'Building placed!'
            mode = null
            await loadVillage()
        } else {
            status.style.color = '#ff6b6b'
            status.textContent = await res.text()
        }
    } else if (mode === 'move') {
        if (!selectedInstanceID) {
            const clicked = buildings.find(b =>
                x >= b.GridX && x < b.GridX + b.Width &&
                y >= b.GridY && y < b.GridY + b.Height
            )
            if (clicked) {
                selectedInstanceID = clicked.InstanceID
                status.style.color = '#6bff6b'
                status.textContent = `Selected ${clicked.Name} — now click where to move it`
            } else {
                status.style.color = '#ff6b6b'
                status.textContent = 'No building found at this location. Click a building first.'
            }
        } else { 
            const res = await fetch('/village/move', {
                method: 'PUT',
                headers: {
                    'Authorization': token,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ instance_id: selectedInstanceID, X: x, Y: y })
            })
            if (res.ok) {
                status.style.color = '#6bff6b'
                status.textContent = 'Building moved!'
            } else {
                status.style.color = '#ff6b6b'
                status.textContent = await res.text()
            }
            
            selectedInstanceID = null
            mode = null
            await loadVillage()
        }   
    }
})

document.getElementById('placeBtn').addEventListener('click', () => {
    mode = 'place'
    status.style.color = '#6bff6b'
    status.textContent = 'Click on grid to place building'
})

document.getElementById('moveBtn').addEventListener('click', () => {
    mode = 'move'
    status.style.color = '#6bff6b'
    status.textContent = 'Click a building to select it'
})

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
    buildings = await response.json()
    //console.log("buildings", buildings)
    render(buildings)
}
loadVillage()
loadShop()