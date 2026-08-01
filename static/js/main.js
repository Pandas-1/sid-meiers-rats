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
let currentOpponent = null
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
    const { x, y } = toGrid(e.clientX - rect.left, e.clientY - rect.top)

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
    } else {
    // no mode — check if clicked a building
    const clicked = buildings.find(b =>
        x >= b.GridX && x < b.GridX + b.Width &&
        y >= b.GridY && y < b.GridY + b.Height
    )
    if (clicked) {
        await showBuildingPopup(clicked)
    }
}
})

async function loadTroops() {
    // fetch both shop troops and current army simultaneously
    const [troopRes, armyRes, levelRes] = await Promise.all([
        fetch('/shop/troops', { headers: { 'Authorization': token } }),
        fetch('/army', { headers: { 'Authorization': token } }),
        fetch('/army/troops', {headers: { 'Authorization': token}})
    ])

    availableTroops = await troopRes.json()
    const currentArmy = await armyRes.json()
    const troopLevels = await levelRes.json()

    // build lookup of current quantities
    const currentComp = {}
    if (currentArmy.ArmyComposition) {
        currentArmy.ArmyComposition.forEach(c => {
            currentComp[c.troop_id] = c.quantity
        })
    }

    // level per troop of the user for the cost formula

    const levelByTroop = {}
    troopLevels.forEach(t => {
        levelByTroop[t.TroopID] = t.TroopLevel
    })

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
        const level = levelByTroop[t.TroopID] || 1
        const row = document.createElement('div')
        row.className = 'troop-row'
        row.innerHTML = `
            <span>${t.Name}</span>
            <span>${t.TroopArmySpace}</span>
            <span>${t.BaseCost*level}</span>
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

    //if (composition.length === 0) {
     //   status.textContent = 'Select at least one troop'
     //   status.style.color = '#ff6b6b'
     //   return
    //}

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
        await loadResources()
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

async function loadTrophies() {
    try {
        const res = await fetch('/village/battlehistory', {
            headers: { 'Authorization': token }
        });
        
        if (res.ok) {
            const battleHistory = await res.json();
            document.getElementById('trophiesDisplay').textContent = `Trophies: ${battleHistory.trophies}`;
        } else {
            const errorText = await res.text();
            console.error("3. Request failed. Server said:", errorText);
        }
    } catch (error) {
        console.error("Network or Fetch error:", error);
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



async function findOpponent() {
    const res = await fetch('/matchmaking/find', {
        method: 'POST',
        headers: { 'Authorization': token }
    })
    if (res.ok) {
        currentOpponent = await res.json()
        console.log("opponent:", currentOpponent)
        document.getElementById('status').textContent = 
            `Found: ${currentOpponent.opponent_username} (${currentOpponent.opponent_trophies} trophies)`
        document.getElementById('attackBtn').style.display = 'block'
    } else {
        status.style.color = '#ff6b6b'
        status.textContent = 'No opponent found'
    }
}

async function startBattle() {
    console.log("startBattle called, currentOpponent:", currentOpponent)
    if (!currentOpponent) {
        console.log("currentOpponent is null!")
        return
    }
    console.log("defender_id being sent:", currentOpponent.opponent_id)
    const res = await fetch('/battle/start', {
        method: 'POST',
        headers: {
            'Authorization': token,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ defender_id: currentOpponent.opponent_id })
    })
    if (res.ok) {
        const data = await res.json()
            window.location.href = `/static/battle.html?battle_id=${data.battle_id}&role=attacker`
    }   else {
        const errorText = await res.text()
            status.style.color = '#ff6b6b'
            status.textContent = errorText  
    }
}

async function loadTroopUpgrades() {
    const [troopRes, userTroopRes] = await Promise.all([
        fetch('/shop/troops', { headers: { 'Authorization': token } }),
        fetch('/army/troops', { headers: { 'Authorization': token } })
    ])

    const allTroops = await troopRes.json()
    const userTroops = await userTroopRes.json()

    // build level lookup by troop_id
    const troopLevels = {}
    if (userTroops) {
        userTroops.forEach(t => {
            troopLevels[t.TroopID] = t.TroopLevel
        })
    }

    const list = document.getElementById('troopUpgradeList')
    list.innerHTML = ''

    allTroops.forEach(t => {
        const level = troopLevels[t.TroopID] || 1
        const upgradeCost = t.BaseCost * level * 2

        const row = document.createElement('div')
        row.className = 'upgrade-row'
        row.innerHTML = `
            <span>${t.Name}</span>
            <span>Lvl ${level}</span>
            <span>Cost: ${upgradeCost} elixir</span>
            <button class="upgrade-btn" data-troop-id="${t.TroopID}">Upgrade</button>
        `
        list.appendChild(row)
    })

    list.querySelectorAll('.upgrade-btn').forEach(btn => {
        btn.addEventListener('click', async () => {
            const troopID = parseInt(btn.dataset.troopId)
            const res = await fetch('/army/upgrade', {
                method: 'PUT',
                headers: {
                    'Authorization': token,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ troop_id: troopID })
            })
            if (res.ok) {
                status.style.color = '#6bff6b'
                status.textContent = 'Troop upgraded!'
                await loadTroopUpgrades()
                await loadResources()
            } else {
                status.style.color = '#ff6b6b'
                status.textContent = await res.text()
            }
        })
    })
}

async function showBuildingPopup(building) {
    const res = await fetch(`/building/${building.InstanceID}/info`, {
        headers: { 'Authorization': token }
    })
    if (!res.ok) return
    const info = await res.json()

    let popup = document.getElementById('buildingPopup')
    if (!popup) {
        popup = document.createElement('div')
        popup.id = 'buildingPopup'
        document.body.appendChild(popup)
    }

    popup.innerHTML = `
        <h3>${building.Name} (Level ${info.current_level})</h3>
        <hr>
        <p><b>Current Stats</b></p>
        <p>HP: ${info.current.HealthBar}</p>
        <p>Attack: ${info.current.DefenceAttack}</p>
        <hr>
        <p><b>Next Level Stats</b></p>
        <p>HP: ${info.next.HealthBar}</p>
        <p>Attack: ${info.next.DefenceAttack}</p>
        <hr>
        <p>Upgrade Cost: ${info.upgrade_cost_r1} gold / ${info.upgrade_cost_r2} elixir</p>
        <div style="display:flex; gap:8px; margin-top:10px">
            <button id="upgradeConfirmBtn">Upgrade</button>
            <button id="closePopupBtn">Close</button>
        </div>
    `
    popup.style.display = 'block'

    document.getElementById('closePopupBtn').addEventListener('click', () => {
        popup.style.display = 'none'
    })

    document.getElementById('upgradeConfirmBtn').addEventListener('click', async () => {
        const res = await fetch('/village/upgrade', {
            method: 'PUT',
            headers: {
                'Authorization': token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ instance_id: building.InstanceID })
        })
        if (res.ok) {
            status.style.color = '#6bff6b'
            status.textContent = 'Building upgraded!'
            popup.style.display = 'none'
            await loadVillage()
            await loadResources()
        } else {
            status.style.color = '#ff6b6b'
            status.textContent = await res.text()
            popup.style.display = 'none'
        }
    })
}

setInterval(async () => {
    await loadResources()
}, 30000)

document.getElementById('matchmakingBtn').addEventListener('click', findOpponent)
document.getElementById('attackBtn').addEventListener('click', startBattle)

loadVillage()
loadShop()
loadResources()
loadTrophies()

document.getElementById('troopUpgradeBtn').addEventListener('click', async () => {
    await loadTroopUpgrades()
    document.getElementById('troopUpgradePanel').style.display = 'block'
})

document.getElementById('closeTroopUpgradeBtn').addEventListener('click', () => {
    document.getElementById('troopUpgradePanel').style.display = 'none'
})
