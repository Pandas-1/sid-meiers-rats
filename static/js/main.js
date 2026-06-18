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
let availableTroops = []
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

async function loadTroops() {
    // fetch both shop troops and current army simultaneously
    const [troopRes, armyRes] = await Promise.all([
        fetch('/shop/troops', { headers: { 'Authorization': token } }),
        fetch('/army', { headers: { 'Authorization': token } })
    ])

    availableTroops = await troopRes.json()
    const currentArmy = await armyRes.json()

    // build lookup of current quantities
    const currentComp = {}
    if (currentArmy.ArmyComposition) {
        currentArmy.ArmyComposition.forEach(c => {
            currentComp[c.troop_id] = c.quantity
        })
    }

    const troopList = document.getElementById('troopList')
    troopList.innerHTML = `
        <div class="troop-row troop-header">
            <span>Troop</span>
            <span>Space</span>
            <span>Cost</span>
            <span>Qty</span>
        </div>
    `

    availableTroops.forEach(t => {
        const currentQty = currentComp[t.TroopID] || 0
        const row = document.createElement('div')
        row.className = 'troop-row'
        row.innerHTML = `
            <span>${t.Name}</span>
            <span>${t.TroopArmySpace}</span>
            <span>${t.BaseCost}</span>
            <input type="number" min="0" value="${currentQty}" id="troop_${t.TroopID}">
        `
        troopList.appendChild(row)
    })

    // show max capacity
    const capacityDiv = document.createElement('div')
    capacityDiv.id = 'capacity'
    capacityDiv.textContent = `Army capacity: ${currentArmy.TroopUnitsUsed} used`
    troopList.appendChild(capacityDiv)
}

async function trainArmy() {
    const composition = []

    availableTroops.forEach(t => {
        const input = document.getElementById(`troop_${t.TroopID}`)
        const qty = parseInt(input.value)
        if (qty > 0) {
            composition.push({ troop_id: t.TroopID, quantity: qty })
        }
    })

    if (composition.length === 0) {
        status.textContent = 'Select at least one troop'
        status.style.color = '#ff6b6b'
        return
    }

    const res = await fetch('/army/train', {
        method: 'POST',
        headers: {
            'Authorization': token,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ composition })
    })

    if (res.ok) {
        status.style.color = '#6bff6b'
        status.textContent = 'Army trained!'
        document.getElementById('armyPanel').style.display = 'none'
    } else {
        status.style.color = '#ff6b6b'
        status.textContent = await res.text()
    }
}

async function loadResources() {
    console.log("loading resources...")
    const res = await fetch('/city', {
        headers: { 'Authorization': token }
    })
    console.log("resources response:", res.status)
    if (res.ok) {
        const city = await res.json()
        document.getElementById('goldDisplay').textContent = `Gold: ${city.Resource1}`
        document.getElementById('elixirDisplay').textContent = `Elixir: ${city.Resource2}`
    }
}

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

document.getElementById('armyBtn').addEventListener('click', async () => {
    await loadTroops()
    document.getElementById('armyPanel').style.display = 'block'
})

document.getElementById('trainBtn').addEventListener('click', trainArmy)

document.getElementById('closeArmyBtn').addEventListener('click', () => {
    document.getElementById('armyPanel').style.display = 'none'
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
loadResources()
