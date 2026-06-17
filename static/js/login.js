document.addEventListener('DOMContentLoaded', function() {
const messsage = document.getElementById('message')

async function register() {
    const username = document.getElementById('username').value
    const password = document.getElementById('password').value

    const response = await fetch('/register', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({username, password})
    })

    if (response.ok) {
        message.textContent = 'Registered! You can now login.'
    } else {
        const text = await response.text()
        message.textContent = text
    }
}

async function login() {
    const username = document.getElementById('username').value
    const password = document.getElementById('password').value

    const response = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    })

    if (response.ok) {
        const data = await response.json()
        localStorage.setItem('token', data.token)
        window.location.href = '/static/index.html'
    } else {
        message.textContent = 'Invalid username or password'
    }
}

document.getElementById('loginBtn').addEventListener('click', login)
document.getElementById('registerBtn').addEventListener('click', register)
})