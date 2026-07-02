const ctx = initCanvas()
drawGrid(ctx)

const token = localStorage.getItem('token')
if (!token) window.location.href = '/static/login.html'

const params = new URLSearchParams(window.location.search)
const battleID = params.get('battle_id')
console.log("battle_id:", battleID)
if (!battleID) window.location.href = '/static/index.html'

let selectedTroopID = null
let army = []
let remainingTroops = {}

async function loadArmy() {
    const res = await fetch('/army', {
        headers: { 'Authorization': token }
    })
    const data = await res.json()
    army = data.ArmyComposition || []

    army.forEach(comp => {
        remainingTroops[comp.troop_id] = comp.quantity
    })

    const troopRes = await fetch('/shop/troops', {
        headers: { 'Authorization': token }
    })
    const allTroops = await troopRes.json()

    const selector = document.getElementById('troopSelector')
    army.forEach(comp => {
        const troop = allTroops.find(t => t.TroopID === comp.troop_id)
        if (!troop) return
        const btn = document.createElement('button')
        btn.className = 'troop-btn'
        btn.textContent = `${troop.Name} x${comp.quantity}`
        btn.dataset.troopId = comp.troop_id
        btn.addEventListener('click', () => {
            document.querySelectorAll('.troop-btn').forEach(b => b.classList.remove('selected'))
            btn.classList.add('selected')
            selectedTroopID = comp.troop_id
            document.getElementById('status').textContent = `Selected: ${troop.Name}`
        })
        selector.appendChild(btn)
    })
}

// connect WebSocket
const ws = new WebSocket(`ws://localhost:8080/battle/ws?battle_id=${battleID}&token=${token}`)

ws.onmessage = (event) => {
    const state = JSON.parse(event.data)

    renderBattle(ctx, state)

    document.getElementById('destruction').textContent = `Destruction: ${state.destruction}%`
    document.getElementById('elapsed').textContent = `Ticks: ${state.elapsed_ticks}`

    if (state.done) {
        ws.close()
        showResult(state)
    }
}

ws.onerror = (e) => {
    document.getElementById('status').textContent = 'Connection error'
}

ctx.canvas.addEventListener('click', (e) => {
    if (!selectedTroopID) {
        document.getElementById('status').textContent = 'Select a troop first'
        return
    }

    // check if troops remaining
    if (!remainingTroops[selectedTroopID] || remainingTroops[selectedTroopID] <= 0) {
        document.getElementById('status').textContent = 'No more troops of this type!'
        return
    }

    const rect = ctx.canvas.getBoundingClientRect()
    const x = Math.floor((e.clientX - rect.left) / CELL_SIZE)
    const y = Math.floor((e.clientY - rect.top) / CELL_SIZE)

    ws.send(JSON.stringify({ troop_id: selectedTroopID, x, y }))

    // decrement remaining
    remainingTroops[selectedTroopID]--

    // update button text
    const btn = document.querySelector(`[data-troop-id="${selectedTroopID}"]`)
    if (btn) {
        const troopName = btn.textContent.split(' x')[0]
        btn.textContent = `${troopName} x${remainingTroops[selectedTroopID]}`
    }
})

function showResult(state) {
    const resultScreen = document.getElementById('resultScreen')
    const resultText = document.getElementById('resultText')
    resultScreen.style.display = 'block'
    resultText.textContent = `Destruction: ${state.destruction}% — ${state.destruction > 50 ? 'Victory!' : 'Defeat'}`
}

document.getElementById('returnBtn').addEventListener('click', () => {
    window.location.href = '/static/index.html'
})

loadArmy()
const role = params.get('role') || 'spectator'
if (role !== 'attacker') {
    document.getElementById('troopSelector').style.display = 'none'
    document.getElementById('status').textContent = 'Spectating battle...'
}