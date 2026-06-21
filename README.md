# sid-meiers-rats

Make sure Docker Desktop is running, then:

```bash
docker compose up --build
```

This will spin up the PostgreSQL database and the Go server automatically.

Visit `http://localhost:8080/static/login.html` to register or login.

## Test Users

Five test users are pre-seeded for testing battles:

| Username | Password |
|----------|----------|
| testuser1 | testuser1 |
| testuser2 | testuser2 |
| testuser3 | testuser3 |
| testuser4 | testuser4 |
| testuser5 | testuser5 |


## Project Structure

```
├── controller/     # HTTP handlers
├── battle/         # Battle simulation and WebSocket logic
├── models/         # Database models and queries
├── routes/         # Route registration
├── db/             # Database connection
├── migrations/     # SQL schema and seed files
└── static/         # Frontend HTML, CSS, JS