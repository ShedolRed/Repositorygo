# Go + MongoDB API (Railway-ready)

## Endpoints
- `GET /health` — checks Mongo connection
- `POST /items` — create item (any JSON)
- `GET /items` — list items
- `GET /items/:id` — get by id
- `PUT /items/:id` — replace `data` with new JSON
- `DELETE /items/:id` — delete by id

## Environment variables
- `PORT` — Railway provides automatically
- `MONGO_URL` — Railway MongoDB service provides automatically
- `MONGO_DB_NAME` — default: `appdb`

## Local run
```bash
export MONGO_URL="mongodb://localhost:27017"
export MONGO_DB_NAME="appdb"
go mod tidy
go run ./cmd/api
```

## Railway deploy (high level)
1) Push this repo to GitHub
2) Railway → New Project → Deploy from GitHub Repo
3) Add → Database → MongoDB
4) In API service → Variables: set `MONGO_DB_NAME=appdb`
5) Generate Domain → open `https://<domain>/health`
