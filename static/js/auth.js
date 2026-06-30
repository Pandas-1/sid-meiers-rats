async function checkAuth() {
    const token = localStorage.getItem('token')
    if (!token) {
        window.location.href = '/static/login.html'
        return
    }

    // hit any protected endpoint to verify token isn't expired
    const res = await fetch('/village', {
        headers: { 'Authorization': token }
    })

    if (res.status === 401) {
        localStorage.removeItem('token')
        window.location.href = '/static/login.html'
    }
}

checkAuth()